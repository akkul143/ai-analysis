import { createInterface } from 'node:readline';
import { handleRequest } from './worker.js';

const lines = createInterface({ input: process.stdin, crlfDelay: Infinity });

for await (const line of lines) {
  const trimmed = line.trim();
  if (trimmed.length === 0) {
    continue;
  }
  let input: unknown;
  try {
    input = JSON.parse(trimmed);
  } catch {
    process.stdout.write(`${JSON.stringify({ id: 'unknown', status: 'error', error: 'invalid_json' })}\n`);
    continue;
  }
  process.stdout.write(`${JSON.stringify(handleRequest(input))}\n`);
}
