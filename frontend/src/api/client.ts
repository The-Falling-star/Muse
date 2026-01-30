import { createClient, ConnectError, Code, type Interceptor } from '@connectrpc/connect';
import {createConnectTransport} from '@connectrpc/connect-web';
import {
    UserService,
    CharacterService,
    ChatService,
    PresetService,
    RegexRuleService,
    WorldInfoService
} from '@/gen/muse/muse_pb';

// 错误处理拦截器
const errorHandlerInterceptor: Interceptor = (next) => async (req) => {
  try {
    return await next(req);
  } catch (error) {
    if (error instanceof ConnectError) {
      const errorMessage = error.message;

        console.error('[API Error]', errorMessage);

      // 特殊处理：未认证时可以跳转登录页
      if (error.code === Code.Unauthenticated) {
        // 可以在这里添加跳转登录页的逻辑
        // router.push('/login');
      }
    } else {
      console.error('网络错误，请检查网络连接', error);
    }

    // 重新抛出错误，让调用方可以选择性处理
    throw error;
  }
};

// 创建 HTTP 传输层（使用 Connect 协议，基于 HTTP + JSON）
const transport = createConnectTransport({
    baseUrl: 'http://localhost:8080',
    interceptors: [errorHandlerInterceptor],
});

// 创建各个服务的客户端
export const userClient = createClient(UserService, transport);
export const characterClient = createClient(CharacterService, transport);
export const chatClient = createClient(ChatService, transport);
export const presetClient = createClient(PresetService, transport);
export const regexRuleClient = createClient(RegexRuleService, transport);
export const worldInfoClient = createClient(WorldInfoService, transport);
