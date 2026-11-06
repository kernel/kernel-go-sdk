import assert from 'node:assert/strict';
import { spawnSync } from 'node:child_process';
import {
  chmodSync,
  copyFileSync,
  mkdirSync,
  mkdtempSync,
  readdirSync,
  readFileSync,
  renameSync,
  rmSync,
  symlinkSync,
  writeFileSync,
} from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import test from 'node:test';

const workflow = readFileSync(new URL('../workflows/update-cli-coverage.yml', import.meta.url), 'utf8');
const jobTexts = Object.fromEntries(
  workflow
    .slice(workflow.indexOf('\njobs:\n'))
    .split(/^  (?=[\w-]+:\n)/m)
    .slice(1)
    .map((job) => [job.slice(0, job.indexOf(':')), job]),
);
const jobs = Object.fromEntries(
  Object.entries(jobTexts).map(([name, job]) => [name, job.split(/(?=^      - )/m).slice(1)]),
);
const steps = Object.values(jobs).flat();
const step = (job, name) => {
  const value = jobs[job].find((entry) => entry.startsWith(`      - name: ${name}\n`));
  assert.ok(value, `missing ${job} step: ${name}`);
  return value;
};
const scriptPath = new URL('../scripts/cli-update-disallowed-paths.sh', import.meta.url).pathname;
const disallowedPaths = (paths) =>
  spawnSync('bash', [scriptPath], {
    input: paths.map((path) => `${path}\0`).join(''),
    encoding: 'utf8',
  })
    .stdout.split('\0')
    .filter(Boolean);

// The run script of a step, as the runner executes it.
const runScript = (entry) => {
  const lines = entry.slice(entry.indexOf('        run: |\n') + '        run: |\n'.length).split('\n');
  return lines.map((line) => line.replace(/^ {10}/, '')).join('\n');
};
const gitEnv = {
  ...process.env,
  GIT_CONFIG_GLOBAL: '/dev/null',
  GIT_CONFIG_NOSYSTEM: '1',
  GIT_AUTHOR_NAME: 'test',
  GIT_AUTHOR_EMAIL: 'test@example.com',
  GIT_COMMITTER_NAME: 'test',
  GIT_COMMITTER_EMAIL: 'test@example.com',
};
const git = (cwd, ...args) => {
  const result = spawnSync('git', args, { cwd, env: gitEnv, encoding: 'utf8' });
  assert.equal(result.status, 0, `git ${args.join(' ')}: ${result.stderr}`);
  return result.stdout.trim();
};
const write = (root, path, content, mode) => {
  mkdirSync(join(root, path, '..'), { recursive: true });
  writeFileSync(join(root, path), content);
  if (mode) chmodSync(join(root, path), mode);
};
const runStep = (entry, cwd, env) =>
  spawnSync('bash', ['-e', '-c', runScript(entry)], {
    cwd,
    env: { ...gitEnv, ...env },
    encoding: 'utf8',
  });

// A temporary workspace with a kernel/cli-like repo on main.
const fixture = (t) => {
  const root = mkdtempSync(join(tmpdir(), 'cli-coverage-'));
  t.after(() => rmSync(root, { recursive: true, force: true }));
  mkdirSync(join(root, '.github/scripts'), { recursive: true });
  copyFileSync(scriptPath, join(root, '.github/scripts/cli-update-disallowed-paths.sh'));
  mkdirSync(join(root, 'updater'));
  symlinkSync(join(root, '.github'), join(root, 'updater/.github'));
  const cli = join(root, 'cli');
  mkdirSync(cli);
  git(cli, 'init', '-q', '-b', 'main');
  write(cli, 'go.mod', 'module example.com/cli\n');
  write(cli, 'cmd/browsers.go', 'package cmd\n');
  write(cli, 'scripts/test', '#!/bin/sh\n', 0o755);
  git(cli, 'add', '-A');
  git(cli, 'commit', '-q', '-m', 'main');
  const runnerTemp = join(root, 'runner-temp');
  mkdirSync(join(runnerTemp, 'cli-coverage-update'), { recursive: true });
  return { root, cli, runnerTemp, mainSha: git(cli, 'rev-parse', 'HEAD'), mainTree: git(cli, 'rev-parse', 'HEAD^{tree}') };
};

// Stages edits on a scratch branch and writes them as the validated patch.
const writePatch = ({ cli, runnerTemp, mainSha }, edit) => {
  git(cli, 'switch', '-q', '-c', 'scratch', mainSha);
  edit(cli);
  git(cli, 'add', '-A');
  writeFileSync(
    join(runnerTemp, 'cli-coverage-update/changes.patch'),
    `${git(cli, 'diff', '--cached', '--binary', '--no-renames')}\n`,
  );
  git(cli, 'reset', '-q', '--hard');
  git(cli, 'clean', '-ffdq');
  git(cli, 'switch', '-q', '--detach', mainSha);
  git(cli, 'branch', '-q', '-D', 'scratch');
};

const applyChanges = (fx, env = {}) =>
  runStep(step('publish', 'Apply validated changes'), fx.cli, {
    RUNNER_TEMP: fx.runnerTemp,
    SDK_VERSION: 'v1.2.3',
    MAIN_SHA: fx.mainSha,
    UPDATE_SHA: '',
    REBASED: '',
    AGENT_BASE_TREE: fx.mainTree,
    ...env,
  });

test('actions are pinned and checkouts do not persist credentials', () => {
  for (const entry of steps.filter((entry) => entry.includes('uses: '))) {
    assert.match(entry, /uses: [\w/-]+@[a-f0-9]{40}(?:\s|$)/);
    if (entry.includes('uses: actions/checkout@')) {
      assert.match(entry, /persist-credentials: false/);
      assert.doesNotMatch(entry, /token:/);
    }
  }
});

test('dispatch inputs reach scripts only through env and are validated', () => {
  // Expressions are substituted into the script before bash parses it.
  for (const entry of steps.filter((entry) => entry.includes('        run: |'))) {
    assert.doesNotMatch(runScript(entry), /\$\{\{[^}]*\binputs\./);
  }
  const prInfo = step('prefetch', 'Get PR info for manual dispatch');
  assert.match(prInfo, /PR_NUMBER_INPUT: \$\{\{ inputs\.pr_number \}\}\n/);
  assert.match(prInfo, /\[\[ ! "\$PR_NUMBER_INPUT" =~ \^\[0-9\]\+\$ \]\][\s\S]*exit 1/);
  assert.ok(jobTexts.generate.includes('needs: prefetch\n'));
});

test('credentialed steps only run in the production SDK repository', () => {
  // The staging mirror shares this file and must not publish unreleased SDK versions.
  assert.match(workflow, /\n  prefetch:\n    if: github\.repository == 'kernel\/kernel-go-sdk'\n/);
  assert.match(workflow, /\n  generate:\n    needs: prefetch\n    if: github\.repository == 'kernel\/kernel-go-sdk'\n/);
});

test('the Pi job references no secrets other than the ones Pi uses', () => {
  // Every secret a job references is loaded into the runner at job start, and
  // Pi could read it from runner memory with passwordless sudo. Pi uses the
  // model API keys and the Kernel API key for its smoke tests. The context key
  // only decrypts the context Pi reads in plaintext anyway.
  const piJobs = Object.keys(jobTexts).filter((name) => /^\s*pi --print\b/m.test(jobTexts[name]));
  assert.deepEqual(piJobs, ['generate']);
  const referenced = (text) => [...text.matchAll(/\bsecrets\b(?:\.([\w-]+))?/g)].map((match) => match[1] ?? '<all>');
  assert.deepEqual(referenced(workflow.slice(0, workflow.indexOf('\njobs:\n'))), []);
  for (const name of piJobs) {
    assert.deepEqual(
      [...new Set(referenced(jobTexts[name]))].sort(),
      ['ANTHROPIC_API_KEY', 'CLI_COVERAGE_CONTEXT_KEY', 'KERNEL_API_KEY', 'OPENAI_API_KEY'],
      `${name} may only reference the secrets Pi uses`,
    );
    assert.deepEqual(
      jobs[name].filter((entry) => entry.includes('secrets.CLI_COVERAGE_CONTEXT_KEY')),
      [step(name, 'Decrypt CLI coverage context')],
    );
    assert.doesNotMatch(jobTexts[name], /github\.token|create-github-app-token|^    permissions:/m);
  }
  assert.match(workflow, /^permissions:\n  contents: read\n\n/m);
});

test('API source access is read-only, revoked, and kept out of the Pi job', () => {
  const readToken = step('prefetch', 'Generate API read token');
  assert.match(readToken, /repositories: kernel\n/);
  assert.match(readToken, /permission-contents: read/);
  assert.doesNotMatch(readToken, /permission-[\w-]+: write/);
  const fetch = step('prefetch', 'Fetch API spec');
  assert.match(fetch, /GH_TOKEN: \$\{\{ steps\.api-read-token\.outputs\.token \}\}/);
  assert.match(fetch, /gh api --method DELETE \/installation\/token\s*$/);
  assert.match(fetch, /-f until="\$1"/);
  assert.doesNotMatch(fetch, /gh repo clone/);
  const tokenUses = steps.filter((entry) => entry.includes('steps.api-read-token.outputs.token'));
  assert.deepEqual(tokenUses, [fetch]);
  assert.doesNotMatch(
    step('generate', 'Clone CLI repo'),
    /GH_TOKEN:|GITHUB_TOKEN:|steps\.[\w-]+\.outputs\.token/,
  );
  // Pi can write to the Go caches, so they must never be saved for other runs.
  assert.match(step('generate', 'Setup Go'), /cache: false/);
});

test('private spec content is only uploaded encrypted, briefly retained, and deleted', () => {
  // This repository is public, so anyone signed in can download its artifacts.
  const uploads = steps.filter((entry) => entry.includes('uses: actions/upload-artifact@'));
  const context = step('prefetch', 'Upload CLI coverage context');
  assert.ok(uploads.includes(context));
  assert.match(context, /path: \$\{\{ runner\.temp \}\}\/cli-coverage-context\.cms\n/);
  assert.match(context, /retention-days: 1\n/);
  // A kept context from a failed attempt must not block "Re-run all jobs".
  assert.match(context, /\n          overwrite: true\n/);
  for (const entry of uploads) {
    const paths = entry.slice(entry.indexOf('path: ')).split(/\n {10}(?=[\w-]+:)/)[0];
    assert.doesNotMatch(paths, /cli-coverage-context(?!\.cms\b)|openapi|stainless|\.patch(?!\b)|\/tmp\//);
  }
  const encrypt = step('prefetch', 'Encrypt CLI coverage context');
  assert.ok(jobs.prefetch.indexOf(step('prefetch', 'Fetch API spec')) < jobs.prefetch.indexOf(encrypt));
  assert.ok(jobs.prefetch.indexOf(encrypt) < jobs.prefetch.indexOf(context));
  assert.match(encrypt, /openssl cms -encrypt -binary -aes-256-gcm/);
  // Only ciphertext and the public certificate cross jobs.
  for (const entry of jobs.generate) assert.doesNotMatch(entry, /GITHUB_ENV|GITHUB_OUTPUT.*KEY/);
  assert.doesNotMatch(jobTexts.prefetch, /outputs:[\s\S]*KEY[\s\S]*steps:/);
  const decrypt = step('generate', 'Decrypt CLI coverage context');
  assert.ok(jobs.generate.indexOf(step('generate', 'Download CLI coverage context')) < jobs.generate.indexOf(decrypt));
  assert.ok(jobs.generate.indexOf(decrypt) < jobs.generate.indexOf(step('generate', 'Update CLI coverage')));
  assert.match(decrypt, /rm "\$key"\n/);
  const cleanup = jobTexts.cleanup;
  assert.match(cleanup, /needs: \[prefetch, generate, publish\]\n/);
  assert.match(cleanup, /permissions:\n      actions: write\n    steps:/);
  assert.match(cleanup, /gh api --method DELETE "repos\/\$GITHUB_REPOSITORY\/actions\/artifacts\/\$ARTIFACT_ID"/);
  assert.doesNotMatch(cleanup, /secrets\.|checkout/);
});

// Evaluates the cleanup job's if: condition for the given upstream job results.
const cleanupRuns = ({ generate, publish, artifactId = '42' }) => {
  const condition = jobTexts.cleanup.match(/^    if: (.+)$/m)[1];
  const js = condition
    .replaceAll('always()', 'true')
    .replaceAll('needs.generate.result', JSON.stringify(generate))
    .replaceAll('needs.publish.result', JSON.stringify(publish))
    .replaceAll('needs.prefetch.outputs.context_artifact_id', JSON.stringify(artifactId))
    .replaceAll(/ (==|!=) /g, ' $1= ')
    .replaceAll("'", '"');
  assert.doesNotMatch(js, /needs\.|\w+\(/, `unhandled expression: ${condition}`);
  // Without a status function, GitHub only runs the job when every need succeeded.
  const implicitSuccess = !/\b(always|success|failure|cancelled)\(\)/.test(condition);
  if (implicitSuccess && (generate !== 'success' || publish !== 'success')) return false;
  return Function(`return ${js};`)();
};

test('the context survives a failed generate so failed jobs can be re-run', () => {
  // "Re-run failed jobs" keeps prefetch's result, so a re-run generate needs the
  // artifact prefetch already uploaded.
  assert.equal(cleanupRuns({ generate: 'failure', publish: 'skipped' }), false);
  assert.equal(cleanupRuns({ generate: 'cancelled', publish: 'skipped' }), false);
  assert.equal(cleanupRuns({ generate: 'skipped', publish: 'skipped', artifactId: '' }), false);
  // After a successful generate nothing downloads it again, publish retries included.
  assert.equal(cleanupRuns({ generate: 'success', publish: 'success' }), true);
  assert.equal(cleanupRuns({ generate: 'success', publish: 'skipped' }), true);
  assert.equal(cleanupRuns({ generate: 'success', publish: 'failure' }), true);
  assert.doesNotMatch(jobTexts.publish, /name: cli-coverage-context\n/);
  assert.match(step('generate', 'Download CLI coverage context'), /name: cli-coverage-context\n/);
});

test('cleanup tolerates a context deleted by an earlier attempt', (t) => {
  const root = mkdtempSync(join(tmpdir(), 'cli-cleanup-'));
  t.after(() => rmSync(root, { recursive: true, force: true }));
  write(root, 'bin/gh', '#!/bin/sh\necho "$GH_STUB_OUTPUT"\nexit "$GH_STUB_STATUS"\n', 0o755);
  const remove = (output, status) =>
    runStep(step('cleanup', 'Delete CLI coverage context'), root, {
      PATH: `${join(root, 'bin')}:${process.env.PATH}`,
      GITHUB_REPOSITORY: 'kernel/kernel-go-sdk',
      ARTIFACT_ID: '42',
      GH_STUB_OUTPUT: output,
      GH_STUB_STATUS: String(status),
    });
  assert.equal(remove('', 0).status, 0);
  const missing = remove('{"message":"Not Found"}gh: Not Found (HTTP 404)', 1);
  assert.equal(missing.status, 0);
  assert.match(missing.stdout, /already deleted/);
  const failed = remove('gh: Server Error (HTTP 500)', 1);
  assert.equal(failed.status, 1);
  assert.match(failed.stderr, /HTTP 500/);
});

test('the context round-trips through encryption and the ciphertext hides it', (t) => {
  const root = mkdtempSync(join(tmpdir(), 'cli-context-'));
  t.after(() => rmSync(root, { recursive: true, force: true }));
  const runnerTemp = join(root, 'runner-temp');
  const keys = join(root, 'updater/.github/keys');
  mkdirSync(keys, { recursive: true });
  const openssl = (...args) => {
    const result = spawnSync('openssl', args, { cwd: root, encoding: 'utf8' });
    assert.equal(result.status, 0, result.stderr);
  };
  openssl('req', '-x509', '-newkey', 'rsa:2048', '-nodes', '-days', '1', '-subj', '/CN=test',
    '-keyout', 'key.pem', '-out', join(keys, 'cli-coverage-context.crt'));
  const secret = 'x-cli-skip: private-endpoint-marker';
  write(runnerTemp, 'cli-coverage-context/api/openapi.yaml', `${secret}\n`);
  write(runnerTemp, 'cli-coverage-context/sdk-diff.patch', 'diff\n');

  const encrypt = runStep(step('prefetch', 'Encrypt CLI coverage context'), root, { RUNNER_TEMP: runnerTemp });
  assert.equal(encrypt.status, 0, encrypt.stdout + encrypt.stderr);
  assert.deepEqual(readdirSync(runnerTemp), ['cli-coverage-context.cms']);
  const ciphertext = readFileSync(join(runnerTemp, 'cli-coverage-context.cms'));
  assert.ok(!ciphertext.includes('private-endpoint-marker'));
  openssl('cms', '-cmsout', '-inform', 'DER', '-in', join(runnerTemp, 'cli-coverage-context.cms'), '-out', 'cms.txt', '-print');
  assert.match(readFileSync(join(root, 'cms.txt'), 'utf8'), /authEnvelopedData/);

  mkdirSync(join(runnerTemp, 'encrypted-context'));
  renameSync(join(runnerTemp, 'cli-coverage-context.cms'), join(runnerTemp, 'encrypted-context/cli-coverage-context.cms'));
  const target = join(root, 'decrypted');
  const decrypt = step('generate', 'Decrypt CLI coverage context').replaceAll('/tmp/cli-coverage-context', target);
  const wrongKey = spawnSync('openssl', ['genpkey', '-algorithm', 'RSA', '-pkeyopt', 'rsa_keygen_bits:2048'], { encoding: 'utf8' }).stdout;
  const failed = runStep(decrypt, root, { RUNNER_TEMP: runnerTemp, CONTEXT_KEY: Buffer.from(wrongKey).toString('base64') });
  assert.notEqual(failed.status, 0);
  const result = runStep(decrypt, root, {
    RUNNER_TEMP: runnerTemp,
    CONTEXT_KEY: readFileSync(join(root, 'key.pem')).toString('base64'),
  });
  assert.equal(result.status, 0, result.stdout + result.stderr);
  assert.equal(readFileSync(join(target, 'api/openapi.yaml'), 'utf8'), `${secret}\n`);
  assert.equal(readFileSync(join(target, 'sdk-diff.patch'), 'utf8'), 'diff\n');
  assert.deepEqual(readdirSync(runnerTemp), []);
});

test('changes are restricted to the allowlist and validated against published modules', () => {
  const agent = step('generate', 'Update CLI coverage');
  const snapshot = step('generate', 'Snapshot publishable changes');
  const verify = step('generate', 'Verify CLI changes');
  assert.ok(jobs.generate.indexOf(agent) < jobs.generate.indexOf(snapshot));
  assert.ok(jobs.generate.indexOf(snapshot) < jobs.generate.indexOf(verify));
  assert.doesNotMatch(agent, /GH_TOKEN:|GITHUB_TOKEN:|steps\.[\w-]+\.outputs\.token/);
  assert.match(agent, /Do not commit, push, open, or merge PRs/);
  // Publishing rebuilds the branch on the recorded main, so Pi compares against it too.
  assert.match(agent, /diff against \$\{\{ steps\.cli-clone\.outputs\.main_sha \}\}/);
  assert.match(agent, /Do not fetch,/);
  assert.doesNotMatch(agent, /origin\/main/);
  assert.match(agent, /cli-scratch/);
  assert.match(snapshot, /updater\/\.github\/scripts\/cli-update-disallowed-paths\.sh/);
  assert.match(step('generate', 'Checkout updater scripts'), /sparse-checkout: \.github\/scripts/);
  assert.match(snapshot, /git reset -q --hard\n\s+git clean -fdq/);
  assert.match(snapshot, /AGENT_BASE: \$\{\{ steps\.cli-clone\.outputs\.agent_base \}\}/);
  assert.match(snapshot, /git reset -q --soft "\$AGENT_BASE"\n\s+git add -A/);
  const prepare = step('generate', 'Clone CLI repo');
  assert.match(prepare, /main_sha=\$\(git rev-parse origin\/main\)/);
  assert.match(prepare, /git config user\.email[\s\S]*git rebase origin\/main; then[\s\S]*else\n\s+echo "rebased=true"/);
  assert.match(prepare, /agent_base=\$\(git rev-parse HEAD\)[\s\S]*agent_base_tree=/);
  // Only a branch without human commits may be regenerated from main.
  assert.match(prepare, /git rebase --abort\n[\s\S]*grep -qvxF '260533166\+kernel-internal\[bot\]@users\.noreply\.github\.com'[\s\S]*exit 1\n[\s\S]*git switch -C /);
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
  // Publishing rebuilds the exact head the changes were validated on, so human
  // commits on the update branch are never squashed or dropped.
  assert.match(apply, /if \[\[ "\$REBASED" == true \]\]; then\n\s+git switch -c \S+ "\$UPDATE_SHA"\n\s+git rebase -q "\$MAIN_SHA"/);
  assert.match(apply, /"\$\(git rev-parse 'HEAD\^\{tree\}'\)" != "\$AGENT_BASE_TREE"[\s\S]*exit 1/);
  assert.ok(apply.indexOf('AGENT_BASE_TREE') < apply.indexOf('git apply --index'));
  const mint = jobs.publish.indexOf(mintStep);
  assert.ok(jobs.publish.indexOf(apply) < mint);
  assert.doesNotMatch(jobs.publish.slice(0, mint).join(''), /GH_TOKEN:/);
  const publish = step('publish', 'Open or update CLI coverage PR');
  assert.match(publish, /GH_TOKEN: \$\{\{ steps\.publish-token\.outputs\.token \}\}/);
  assert.match(publish, /gh auth setup-git\n\s+git push --force-with-lease="cli-coverage-update:\$UPDATE_SHA"/);
  for (const entry of jobs.publish) {
    assert.doesNotMatch(entry, /continue-on-error: true|always\(\)/);
  }
});

test('the publish gate commits allowlisted regular files', (t) => {
  const fx = fixture(t);
  writePatch(fx, (cli) => {
    write(cli, 'go.mod', 'module example.com/cli\n\nrequire example.com/sdk v1.2.3\n');
    write(cli, 'cmd/browsers.go', 'package cmd\n\nvar size int\n');
    write(cli, 'pkg/util/client.go', 'package util\n');
    write(cli, 'cmd/kernel/main.go', 'package main\n', 0o755);
  });
  const result = applyChanges(fx);
  assert.equal(result.status, 0, result.stdout + result.stderr);
  assert.equal(git(fx.cli, 'rev-parse', 'HEAD~1'), fx.mainSha);
  assert.deepEqual(git(fx.cli, 'diff', '--name-only', 'HEAD~1', 'HEAD').split('\n'), [
    'cmd/browsers.go',
    'cmd/kernel/main.go',
    'go.mod',
    'pkg/util/client.go',
  ]);
  assert.equal(git(fx.cli, 'log', '-1', '--format=%s'), 'Update Go SDK to v1.2.3 and CLI coverage');
});

test('the publish gate rejects paths outside the allowlist', (t) => {
  const fx = fixture(t);
  writePatch(fx, (cli) => {
    write(cli, 'cmd/browsers.go', 'package cmd\n\nvar size int\n');
    write(cli, 'scripts/test', '#!/bin/sh\ncurl example.com\n', 0o755);
  });
  const result = applyChanges(fx);
  assert.equal(result.status, 1);
  assert.match(result.stdout, /outside the CLI updater allowlist:\nscripts\/test/);
  assert.equal(git(fx.cli, 'rev-parse', 'HEAD'), fx.mainSha);
});

test('the publish gate rejects symlinks and submodules on allowlisted paths', (t) => {
  for (const edit of [
    (cli) => symlinkSync('/etc/passwd', join(cli, 'cmd/passwd.go')),
    (cli) => {
      const nested = join(cli, 'cmd/module.go');
      mkdirSync(nested);
      git(nested, 'init', '-q');
      git(nested, 'commit', '-q', '--allow-empty', '-m', 'nested');
    },
  ]) {
    const fx = fixture(t);
    writePatch(fx, edit);
    const result = applyChanges(fx);
    assert.equal(result.status, 1, result.stdout + result.stderr);
    assert.match(result.stdout, /may only contain regular files/);
    assert.equal(git(fx.cli, 'rev-parse', 'HEAD'), fx.mainSha);
  }
});

test('the publish gate rejects a base that differs from the validated one', (t) => {
  const fx = fixture(t);
  writePatch(fx, (cli) => write(cli, 'cmd/browsers.go', 'package cmd\n\nvar size int\n'));
  const result = applyChanges(fx, { AGENT_BASE_TREE: 'b'.repeat(40) });
  assert.equal(result.status, 1);
  assert.match(result.stdout, /differs from the one the changes were validated on/);
});

test('the snapshot keeps agent commits and drops everything outside the allowlist', (t) => {
  const fx = fixture(t);
  git(fx.cli, 'switch', '-q', '-c', 'cli-coverage-update');
  write(fx.cli, 'cmd/browsers.go', 'package cmd\n\nvar size int\n');
  git(fx.cli, 'commit', '-q', '-am', 'agent commit');
  write(fx.cli, 'pkg/util/client.go', 'package util\n');
  write(fx.cli, 'notes.md', 'scratch\n');
  write(fx.cli, 'scripts/test', '#!/bin/sh\ncurl example.com\n', 0o755);
  const output = join(fx.runnerTemp, 'output');
  const result = runStep(step('generate', 'Snapshot publishable changes'), fx.cli, {
    AGENT_BASE: fx.mainSha,
    GITHUB_WORKSPACE: fx.root,
    GITHUB_OUTPUT: output,
    RUNNER_TEMP: fx.runnerTemp,
  });
  assert.equal(result.status, 0, result.stdout + result.stderr);
  assert.equal(readFileSync(output, 'utf8'), 'changed=true\n');
  assert.equal(git(fx.cli, 'rev-parse', 'HEAD'), fx.mainSha);
  // The validated checkout holds only the allowlisted changes, all staged.
  assert.deepEqual(git(fx.cli, 'status', '--porcelain', '--untracked-files=all').split('\n'), [
    'M  cmd/browsers.go',
    'A  pkg/util/client.go',
  ]);
  const patch = readFileSync(join(fx.runnerTemp, 'changes.patch'), 'utf8');
  assert.doesNotMatch(patch, /notes\.md|scripts\/test/);
});
