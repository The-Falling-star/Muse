import { createClient, ConnectError, Code, type Interceptor } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { globalMessage } from '@/composables/useGlobalMessage';
import {UserService} from "@/gen/muse/user_pb.ts";
import {CharacterService} from "@/gen/muse/character_pb.ts";
import {ChatService} from "@/gen/muse/chat_pb.ts";
import {PresetService} from "@/gen/muse/preset_pb.ts";
import {RegexRuleService} from "@/gen/muse/regex_pb.ts";
import {WorldInfoService} from "@/gen/muse/worldinfo_pb.ts";
import {FileService} from "@/gen/muse/file_pb.ts";
import {CommonService} from "@/gen/muse/common_pb.ts";

// 认证拦截器 - 自动添加token到请求头
const authInterceptor: Interceptor = (next) => async (req) => {
  const token = localStorage.getItem('token');
  if (token) {
    req.header.set('Authorization', `Bearer ${token}`);
  }
  return await next(req);
};

// 错误处理拦截器
const errorHandlerInterceptor: Interceptor = (next) => async (req) => {
  try {
    return await next(req);
  } catch (error) {
    console.error('[Error Interceptor] Caught error:', error);
    console.error('[Error Interceptor] Error type:', error?.constructor?.name);
    console.error('[Error Interceptor] Is ConnectError?', error instanceof ConnectError);

    if (error instanceof ConnectError) {
      const errorMessage = error.message;
      const errorCode = error.code;

      console.error('[API Error]', {
        message: errorMessage,
        code: errorCode,
        codeName: Code[errorCode],
        rawCode: error.code
      });

      // 显示错误消息（除了未认证错误，因为会自动跳转）
      if (errorCode !== Code.Unauthenticated) {
        console.error('[Error Interceptor] Showing error message:', errorMessage);
        globalMessage.error(errorMessage);
      } else {
        console.error('[Error Interceptor] Skipping error message (Unauthenticated)');
      }

      // 特殊处理：未认证时跳转登录页
      if (errorCode === Code.Unauthenticated) {
        // 清除本地token
        localStorage.removeItem('token');
        // 跳转到登录页（避免在登录页循环跳转）
        const loginPath = import.meta.env.BASE_URL + 'login';
        if (!window.location.pathname.endsWith('/login')) {
          window.location.href = loginPath;
        }
      }
    } else {
      const errorMsg = '网络错误，请检查网络连接';
      console.error(errorMsg, error);
      globalMessage.error(errorMsg);
    }

    // 重新抛出错误，让调用方可以选择性处理
    throw error;
  }
};

// 创建 HTTP 传输层（使用 Connect 协议，基于 HTTP + JSON）
const transport = createConnectTransport({
  baseUrl: import.meta.env.VITE_API_BASE_URL,
  interceptors: [authInterceptor, errorHandlerInterceptor],
});

// 创建各个服务的客户端
export const commonClient = createClient(CommonService, transport);
export const userClient = createClient(UserService, transport);
export const characterClient = createClient(CharacterService, transport);
export const chatClient = createClient(ChatService, transport);
export const presetClient = createClient(PresetService, transport);
export const regexRuleClient = createClient(RegexRuleService, transport);
export const worldInfoClient = createClient(WorldInfoService, transport);
export const fileClient = createClient(FileService, transport);
