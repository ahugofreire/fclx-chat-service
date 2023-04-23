import { chatClient } from './client';
import { ChatServiceClient as GrpcChatServiceClient} from './rpc/pb/ChatService';
import {Metadata} from '@grpc/grpc-js';

interface ChatStraemDTO {
    chat_id?: string | null; 
    user_id: string; 
    message: string | null;
}

export class ChatServiceClient {
    private authorization = "123456";

    constructor(private chatclient: GrpcChatServiceClient) {}

    chatStream({ chat_id, user_id, message}: ChatStraemDTO) {
        
        const metadata = new Metadata();
        metadata.set("authorization",this.authorization)

        const stream = this.chatclient.chatStream({
            chatId: chat_id!,
            userId: user_id,
            userMessage: message!,
        }, metadata);

        // stream.on("data", (data) => {
        //     console.log(data)
        // });

        // stream.on("error", (err) => {
        //     console.log(err);
        // });

        // stream.on("end", () => {
        //     console.log("end");
        // });

        return stream;
    }
 }

 export class ChatServiceClientFactory {
    static create() {
        return new ChatServiceClient(chatClient)
    }
 }