import { Injectable } from '@nestjs/common';
import { User } from '@prisma/client';

@Injectable()
export class UserService {
  getMe(user: User) {
    delete user.hash;

    return user;
  }
}
