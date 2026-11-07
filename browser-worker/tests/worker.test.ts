import assert from 'node:assert/strict';
import test from 'node:test';
import { handleRequest } from '../src/worker.js';

test('health describes the Phase 1 browser boundary', () => {
  const response = handleRequest({ id: 'test-health', action: 'health' });
  assert.equal(response.id, 'test-health');
  assert.equal(response.status, 'ok');
  assert.equal(response.result?.['worker'], 'ai-analysis-browser-worker');
  assert.equal(response.result?.['phase'], 1);
});

test('collection is explicitly deferred to Phase 2', () => {
  const response = handleRequest({ id: 'test-collect', action: 'collect' });
  assert.equal(response.status, 'not_implemented');
  assert.equal(response.error, 'collect_is_added_in_phase_2');
});

test('invalid requests fail closed', () => {
  const response = handleRequest({ action: 'health' });
  assert.equal(response.status, 'error');
  assert.equal(response.error, 'invalid_request');
});
