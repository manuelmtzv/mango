import { createParamDecorator, ExecutionContext } from '@nestjs/common';

export const GetJwt = createParamDecorator(
  (data: string, ctx: ExecutionContext) => {
    const request = ctx.switchToHttp().getRequest();
    return extractPlainToken(request.headers.authorization);
  },
);

function extractPlainToken(token: string) {
  return token.split(' ')[1];
}
