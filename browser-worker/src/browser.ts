import { existsSync } from 'node:fs';
import { chromium } from 'playwright';

export const browserCapability = (): Record<string, unknown> => {
  const executablePath = chromium.executablePath();
  return {
    engine: 'chromium',
    executableInstalled: existsSync(executablePath)
  };
};
