import assert from 'node:assert/strict';
import { spawnSync } from 'node:child_process';
import { readFileSync } from 'node:fs';
import test from 'node:test';

const workflow = readFileSync(new URL('../workflows/update-cli-coverage.yml', import.meta.url), 'utf8');
const jobs = Object.fromEntries(
  workflow
    .slice(workflow.indexOf('\njobs:\n'))
    .split(/^  (?=[\w-]+:\n)/m)
    .slice(1)
    .map((job) => [job.slice(0, job.indexOf(':')), job.split(/(?=^      - )/m).slice(1)]),
);
const steps = Object.values(jobs).flat();
const step = (job, name) => {
  const value = jobs[job].find((entry) => entry.startsWith(`      - name: ${name}\n`));
  assert.ok(value, `missing ${job} step: ${name}`);
  return value;
};
const disallowedPaths = (paths) =>
  spawnSync('bash', [new URL('../scripts/cli-update-disallowed-paths.sh', import.meta.url).pathname], {
    input: paths.map((path) => `${path}\0`).join(''),
    encoding: 'utf8',
  })
    .stdout.split('\0')
    .filter(Boolean);

test('actions are pinned and checkouts do not persist credentials', () => {
  for (const entry of steps.filter((entry) => entry.includes('uses: '))) {
    assert.match(entry, /uses: [\w/-]+@[a-f0-9]{40}(?:\s|$)/);
    if (entry.includes('uses: actions/checkout@')) {
      assert.match(entry, /persist-credentials: false/);
      assert.doesNotMatch(entry, /token:/);
    }
  }
});

test('API source access is read-only and revoked before generation', () => {
  const readToken = step('generate', 'Generate API read token');
  assert.match(readToken, /repositories: kernel\n/);
  assert.match(readToken, /permission-contents: read/);
  assert.doesNotMatch(readToken, /permission-[\w-]+: write/);
  const fetch = step('generate', 'Fetch API spec');
  assert.match(fetch, /GH_TOKEN: \$\{\{ steps\.api-read-token\.outputs\.token \}\}/);
  assert.match(fetch, /gh api --method DELETE \/installation\/token\s*$/);
  assert.match(fetch, /-f until="\$1"/);
  assert.doesNotMatch(fetch, /gh repo clone/);
  assert.ok(jobs.generate.indexOf(fetch) < jobs.generate.indexOf(step('generate', 'Update CLI coverage')));
  assert.doesNotMatch(
    step('generate', 'Clone CLI repo'),
    /GH_TOKEN:|GITHUB_TOKEN:|steps\.[\w-]+\.outputs\.token/,
  );
});

test('changes are restricted to the allowlist and validated against published modules', () => {
  const agent = step('generate', 'Update CLI coverage');
  const snapshot = step('generate', 'Snapshot publishable changes');
  const verify = step('generate', 'Verify CLI changes');
  assert.ok(jobs.generate.indexOf(agent) < jobs.generate.indexOf(snapshot));
  assert.ok(jobs.generate.indexOf(snapshot) < jobs.generate.indexOf(verify));
  assert.doesNotMatch(agent, /GH_TOKEN:|GITHUB_TOKEN:|steps\.[\w-]+\.outputs\.token/);
  assert.match(agent, /Do not commit, push, open, or merge PRs/);
  assert.match(agent, /cli-scratch/);
  assert.match(snapshot, /updater\/\.github\/scripts\/cli-update-disallowed-paths\.sh/);
  assert.match(step('generate', 'Checkout updater scripts'), /sparse-checkout: \.github\/scripts/);
  assert.match(snapshot, /git reset -q --hard\n\s+git clean -fdq/);
  assert.match(snapshot, /AGENT_BASE: \$\{\{ steps\.cli-clone\.outputs\.agent_base \}\}/);
  assert.match(snapshot, /git reset -q --soft "\$AGENT_BASE"\n\s+git add -A/);
  assert.match(step('generate', 'Clone CLI repo'), /agent_base=\$\(git rev-parse HEAD\)/);
  const scriptCalls = workflow.match(/^.*disallowed-paths\.sh.*$/gm);
  assert.equal(scriptCalls.length, 2);
  for (const line of scriptCalls) assert.match(line, /\| bash /);
  assert.match(verify, /\{\{with \.Replace\}\}/);
  assert.match(verify, /select\(\.New\.Version == null\)/);
  for (const command of ['go build ./...', 'go vet ./...', 'go test ./...']) {
    assert.ok(verify.includes(command), `missing ${command}`);
  }
  assert.deepEqual(
    disallowedPaths([
      'go.mod',
      'go.sum',
      'cmd/browsers.go',
      'cmd/kernel/main.go',
      'pkg/util/client.go',
      'bin/kernel',
      'response.json',
      'cmd/output.json',
      '.github/workflows/release.yaml',
      'pkg/templates/python/sample-app/main.py',
    ]),
    [
      'bin/kernel',
      'response.json',
      'cmd/output.json',
      '.github/workflows/release.yaml',
      'pkg/templates/python/sample-app/main.py',
    ],
  );
});

test('publishing runs in a separate job that never executes repository code', () => {
  assert.doesNotMatch(jobs.generate.join(''), /permission-[\w-]+: write/);
  const mintStep = step('publish', 'Generate CLI publishing token');
  assert.match(mintStep, /repositories: cli\n/);
  assert.match(mintStep, /permission-contents: write/);
  const publishRuns = jobs.publish.filter((entry) => entry.includes('run: |')).join('');
  assert.doesNotMatch(publishRuns, /^\s*(go|bun|npm|npx|node|pi|make)\b/m);
  assert.doesNotMatch(jobs.publish.join(''), /secrets\.(ANTHROPIC|OPENAI|KERNEL)_API_KEY/);
  const apply = step('publish', 'Apply validated changes');
  assert.match(apply, /git apply --index "\$RUNNER_TEMP\/cli-coverage-update\/changes\.patch"/);
  assert.match(apply, /cli-update-disallowed-paths\.sh/);
  assert.match(apply, /100644\|100755\|000000/);
  const mint = jobs.publish.indexOf(mintStep);
  assert.ok(jobs.publish.indexOf(apply) < mint);
  assert.doesNotMatch(jobs.publish.slice(0, mint).join(''), /GH_TOKEN:/);
  const publish = step('publish', 'Open or update CLI coverage PR');
  assert.match(publish, /GH_TOKEN: \$\{\{ steps\.publish-token\.outputs\.token \}\}/);
  assert.match(publish, /gh auth setup-git\n\s+git push/);
  for (const entry of jobs.publish) {
    assert.doesNotMatch(entry, /continue-on-error: true|always\(\)/);
  }
});
