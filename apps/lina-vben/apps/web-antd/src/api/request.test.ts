import { beforeEach, describe, expect, it, vi } from 'vitest';

const { appConfig } = vi.hoisted(() => ({
  appConfig: {
    apiURL: undefined as string | undefined,
  },
}));

vi.mock('@vben/hooks', () => ({
  useAppConfig: () => appConfig,
}));

vi.mock('@vben/preferences', () => ({
  preferences: {
    app: {
      enableRefreshToken: false,
      locale: 'zh-CN',
      loginExpiredMode: 'modal',
    },
  },
}));

vi.mock('@vben/request', () => {
  class MockRequestClient {
    public instance = {
      post: vi.fn(),
    };

    public addRequestInterceptor = vi.fn();

    public addResponseInterceptor = vi.fn();

    public constructor(public readonly options?: unknown) {}
  }

  return {
    authenticateResponseInterceptor: vi.fn(() => vi.fn()),
    defaultResponseInterceptor: vi.fn(() => vi.fn()),
    errorMessageResponseInterceptor: vi.fn(() => vi.fn()),
    RequestClient: MockRequestClient,
  };
});

vi.mock('@vben/stores', () => ({
  useAccessStore: () => ({
    accessToken: null,
    isAccessChecked: false,
    refreshToken: null,
    setAccessToken: vi.fn(),
    setLoginExpired: vi.fn(),
    setRefreshToken: vi.fn(),
  }),
}));

vi.mock('ant-design-vue', () => ({
  message: {
    error: vi.fn(),
  },
}));

vi.mock('#/locales', () => ({
  $t: (key: string) => key,
}));

vi.mock('#/store', () => ({
  useAuthStore: () => ({
    clearSession: vi.fn(),
  }),
}));

vi.mock('#/store/tenant', () => ({
  useTenantStore: () => ({
    currentTenant: null,
    enabled: false,
  }),
}));

async function importRequestWithApiURL(apiURL?: string) {
  appConfig.apiURL = apiURL;
  vi.resetModules();
  return await import('./request');
}

describe('request API helpers', () => {
  beforeEach(() => {
    appConfig.apiURL = undefined;
  });

  it('uses the default host API base when the test environment omits VITE_GLOB_API_URL', async () => {
    const { pluginApiPath } = await importRequestWithApiURL();

    expect(pluginApiPath('linapro-tenant-core', '/auth/login-tenants')).toBe(
      `${window.location.origin}/x/linapro-tenant-core/api/v1/auth/login-tenants`,
    );
  });

  it('derives plugin API paths from an absolute host API base', async () => {
    const { pluginApiPath } = await importRequestWithApiURL(
      'https://example.test/api/v1/',
    );

    expect(pluginApiPath('/linapro-tenant-core/', 'platform/tenants')).toBe(
      'https://example.test/x/linapro-tenant-core/api/v1/platform/tenants',
    );
  });
});
