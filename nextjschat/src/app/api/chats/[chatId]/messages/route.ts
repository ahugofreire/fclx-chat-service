import { withAuth } from "@/app/api/helpers";
import { prisma } from "@/app/prisma/prisma";
import { NextRequest, NextResponse } from "next/server";

export const GET = withAuth( async (_request: NextRequest,
  token, 
  { params }: { params: { chatId: string}}
 ) => {

  const chat = await prisma.chat.findUniqueOrThrow({
    where: {
        id: params.chatId
    }
  });

  if (chat.user_id !== token.sub) {
    return NextResponse.json({ error: "Not found"}, { status: 404});
  }

    const messages = await prisma.message.findMany({
        where: {
            chat_id: params.chatId
        },
        orderBy: { created_at: "asc" },
    });

    return NextResponse.json(messages)
});

export const POST = withAuth( async (
    request: NextRequest,
    token,
    { params }: { params: { chatId: string }}
) => {
  
  const chat = await prisma.chat.findUniqueOrThrow({
    where: {
        id: params.chatId
    }
  });

  if (chat.user_id !== token.sub) {
    return NextResponse.json({ error: "Not found"}, { status: 404});
  }


  const body = await request.json();
  const messageCreated = await prisma.message.create({
    data: {
        content: body.message,
        chat_id: chat.id
    }
  });

  return NextResponse.json(messageCreated);
});