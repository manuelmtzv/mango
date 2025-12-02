import {
  BadRequestException,
  Injectable,
  NotFoundException,
} from '@nestjs/common';
import { PrismaService } from '../prisma/prisma.service';
import { CreateTagDto } from './dto/CreateTagDto';
import { UpdateTagDto } from './dto/UpdateTag';
import { uniq } from 'rambda';

@Injectable()
export class TagService {
  constructor(private prisma: PrismaService) {}

  async getTags(userId: string) {
    const tags = await this.prisma.tag.findMany({
      where: { userId },
      orderBy: { updatedAt: 'desc' },
    });

    return {
      data: tags,
      count: tags.length,
    };
  }

  async getTag(userId: string, tagId: string) {
    const tag = await this.prisma.tag.findUnique({
      where: { id: tagId, userId },
    });

    if (!tag) throw new NotFoundException('Tag with provided id was not found');
    return tag;
  }

  async findTagsOrCreate(userId: string, tags: string[]) {
    tags = uniq(tags);

    const previouslyCreatedTags = await this.prisma.tag.findMany({
      where: { userId, name: { in: tags } },
    });

    const toCreateTags = tags.filter(
      (tag) => !previouslyCreatedTags.find((t) => t.name === tag),
    );

    if (toCreateTags.length === 0) {
      return {
        data: previouslyCreatedTags,
        count: previouslyCreatedTags.length,
      };
    }

    if (previouslyCreatedTags.length + toCreateTags.length > 50) {
      throw new BadRequestException('You can only have a maximum of 50 tags.');
    }

    await this.prisma.tag.createMany({
      data: toCreateTags.map((tag) => ({ name: tag, userId })),
    });

    const allTags = await this.prisma.tag.findMany({
      where: { userId, name: { in: tags } },
    });

    return {
      data: allTags,
      count: allTags.length,
    };
  }

  async createTag(userId: string, dto: CreateTagDto) {
    const userTags = await this.prisma.tag.findMany({ where: { userId } });

    if (userTags.length >= 50) {
      throw new BadRequestException('You can only have 50 tags.');
    }

    for (const tag of userTags) {
      if (tag.name === dto.name) {
        throw new BadRequestException(
          `Tag with provided name (${dto.name}) already exists.`,
        );
      }
    }

    const tag = this.prisma.tag.create({
      data: {
        ...dto,
        userId,
      },
    });

    return tag;
  }

  async updateTag(userId: string, tagId: string, dto: UpdateTagDto) {
    const tag = await this.prisma.tag.findUnique({
      where: { id: tagId, userId },
    });

    if (!tag)
      throw new NotFoundException('Tag with provided id was not found.');

    const updatedTag = this.prisma.tag.update({
      where: { id: tagId, userId },
      data: dto,
    });

    return updatedTag;
  }

  async deleteTag(userId: string, tagId: string) {
    const tag = await this.prisma.tag.findUnique({
      where: { id: tagId, userId },
    });

    if (!tag)
      throw new NotFoundException('Tag with provided id was not found.');

    await this.prisma.tag.delete({ where: { id: tagId } });
  }
}
