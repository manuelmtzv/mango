import { Controller, Get, UseGuards } from '@nestjs/common';
import { JwtGuard } from '@modules/auth/guard';
import { ApiBearerAuth, ApiTags } from '@nestjs/swagger';
import { GetUser } from '../auth/decorator';
import { User } from '@prisma/client';
import { UserService } from './user.service';

@ApiTags('Users')
@ApiBearerAuth('JWT-auth')
@UseGuards(JwtGuard)
@Controller('users')
export class UserController {
  constructor(private userService: UserService) {}

  @Get('me')
  getMe(@GetUser() user: User) {
    return this.userService.getMe(user);
  }
}
