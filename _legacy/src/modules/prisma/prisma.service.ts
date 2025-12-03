import { Injectable } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { PrismaClient } from '@prisma/client';

@Injectable()
export class PrismaService extends PrismaClient {
  constructor(config: ConfigService) {
    super({
      datasources: {
        db: {
          url: config.get<string>('DATABASE_URL'),
        },
      },
    });
  }

  async cleanDb() {
    return this.$transaction([
      this.tag.deleteMany(),
      this.note.deleteMany(),
      this.user.deleteMany(),
    ]);
  }
}
