import { browserCapability } from './browser.js';
import { decodeRequest, type WorkerResponse } from './protocol.js';

export const handleRequest = (input: unknown): WorkerResponse => {
  const request = decodeRequest(input);
  if (request === null) {
    return {
      id: 'unknown',
      status: 'error',
      error: 'invalid_request'
    };
  }

  if (request.action === 'health') {
    return {
      id: request.id,
      status: 'ok',
      result: {
        worker: 'ai-analysis-browser-worker',
        version: '0.1.0',
        phase: 1,
        ...browserCapability()
      }
    };
  }

  return {
    id: request.id,
    status: 'not_implemented',
    error: `${request.action}_is_added_in_phase_2`
  };
};
