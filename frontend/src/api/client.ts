import {
    UserService,
    ChatService,
    CharacterService,
    PresetService,
    RegexRuleService,
    WorldInfoService,
} from "@/gen/muse/muse_pb.ts";

import {createConnectTransport} from "@connectrpc/connect-web";
import {createClient} from "@connectrpc/connect";


const transport = createConnectTransport({
    baseUrl: import.meta.env.VITE_API_BASE_URL,
});

export const userService = createClient(UserService, transport);
export const chatService = createClient(ChatService, transport);
export const characterService = createClient(CharacterService, transport);
export const presetService = createClient(PresetService, transport);
export const regexRuleService = createClient(RegexRuleService, transport);
export const worldInfoService = createClient(WorldInfoService, transport);