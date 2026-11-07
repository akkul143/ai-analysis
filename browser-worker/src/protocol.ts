export type WorkerAction = 'health' | 'auth_login' | 'collect';

export type WorkerRequest = {
  id: string;
  action: WorkerAction;
  source?: string;
  profileDirectory?: string;
};

export type WorkerResponse = {
  id: string;
  status: 'ok' | 'error' | 'not_implemented';
  result?: Record<string, unknown>;
  error?: string;
};

export const decodeRequest = (input: unknown): WorkerRequest | null => {
  if (typeof input !== 'object' || input === null) {
    return null;
  }
  const candidate = input as Record<string, unknown>;
  const id = candidate['id'];
  const action = candidate['action'];
  if (
    typeof id !== 'string' ||
    !['health', 'auth_login', 'collect'].includes(typeof action === 'string' ? action : '')
  ) {
    return null;
  }
  return {
    id,
    action: action as WorkerAction,
    ...(typeof candidate['source'] === 'string' ? { source: candidate['source'] } : {}),
    ...(typeof candidate['profileDirectory'] === 'string'
      ? { profileDirectory: candidate['profileDirectory'] }
      : {})
  };
};
