import {
  BadRequestException,
  ForbiddenException,
  Injectable,
} from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import { ConfigService } from '@nestjs/config';
import * as argon from 'argon2';
import { omit } from 'rambda';
import { RegisterDto } from './dto/register.dto';
import { LoginDto } from './dto/login.dto';
import { PrismaService } from '@modules/prisma/prisma.service';
import { User } from '@prisma/client';

@Injectable()
export class AuthService {
  constructor(
    private jwt: JwtService,
    private config: ConfigService,
    private prisma: PrismaService,
  ) {}

  async validateToken(user: User, token: string) {
    return { username: user.username, accessToken: token };
  }

  async register(dto: RegisterDto) {
    if (await this.prisma.user.findFirst({ where: { email: dto.email } }))
      throw new BadRequestException('Email already exists');

    if (await this.prisma.user.findFirst({ where: { username: dto.username } }))
      throw new BadRequestException('Username already exists');

    const hash = await argon.hash(dto.password);

    try {
      const user = await this.prisma.user.create({
        data: {
          ...omit(['password'], dto),
          hash,
        },
      });

      delete user.hash;
      return this.signToken(user.id, user.email, user.username);
    } catch (err) {
      throw err;
    }
  }

  async login(dto: LoginDto) {
    const user = await this.prisma.user.findFirst({
      where: {
        OR: [{ email: dto.identifier }, { username: dto.identifier }],
      },
    });

    if (!user) throw new ForbiddenException('Invalid identifier or password');

    const passwordValid = await argon.verify(user.hash, dto.password);

    if (!passwordValid)
      throw new ForbiddenException('Invalid identifier or password');

    return this.signToken(user.id, user.email, user.username);
  }

  async signToken(
    userId: string | number,
    email: string,
    username: string,
  ): Promise<{ accessToken: string; username: string }> {
    const payload = {
      sub: userId,
      email,
    };

    const token = await this.jwt.signAsync(payload, {
      expiresIn: '7d',
      secret: this.config.get('JWT_SECRET'),
    });

    return {
      accessToken: token,
      username,
    };
  }
}
