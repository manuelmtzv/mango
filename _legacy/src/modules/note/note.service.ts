import { Injectable, NotFoundException } from '@nestjs/common';
import { PrismaService } from '@modules/prisma/prisma.service';
import { CreateNoteDto, PaginateNotesDto } from './dto';
import { UpdateNoteDto } from './dto/updateNote.dto';
import { paginate } from '@src/shared/utils/paginate';
import { TagService } from '../tag/tag.service';
import { Tag } from '@prisma/client';

@Injectable()
export class NoteService {
  constructor(
    private prisma: PrismaService,
    private readonly tagService: TagService,
  ) {}

  async getNotes(userId: string, dto: PaginateNotesDto) {
    const { page = 1, limit = 10, search, tags } = dto;

    const noteCount = (
      await this.prisma.note.findMany({
        where: { userId },
      })
    ).length;

    const pagination = paginate({ count: noteCount, page, limit });

    const notes = await this.prisma.note.findMany({
      where: {
        userId,
        title: search && { contains: search },
        tags: tags && {
          some: {
            name: { in: tags },
          },
        },
      },
      orderBy: { updatedAt: 'desc' },
      include: {
        tags: {
          select: {
            id: true,
            name: true,
            color: true,
          },
        },
      },
      ...pagination.config(),
    });

    return {
      data: notes,
      meta: {
        page,
        limit,
        count: notes.length,
        totalPages: pagination.maxPage,
      },
    };
  }

  async getNote(userId: string, noteId: string) {
    const note = await this.prisma.note.findUnique({
      where: { id: noteId, userId },
      include: {
        tags: {
          select: {
            id: true,
            name: true,
            color: true,
          },
        },
      },
    });

    if (!note)
      throw new NotFoundException('Note with provided id was not found');
    return note;
  }

  async createNote(userId: string, dto: CreateNoteDto) {
    let createdTags: Tag[] = [];
    const { tags, ...rest } = dto;

    if (tags) {
      const { data } = await this.tagService.findTagsOrCreate(userId, tags);
      createdTags = data;
    }

    const note = this.prisma.note.create({
      data: {
        ...rest,
        userId,
        tags: {
          connect: createdTags.length
            ? createdTags.map((tag) => ({ id: tag.id }))
            : undefined,
        },
      },
      include: {
        tags: {
          select: {
            id: true,
            name: true,
            color: true,
          },
        },
      },
    });

    return note;
  }

  async updateNote(userId: string, noteId: string, dto: UpdateNoteDto) {
    const note = await this.prisma.note.findUnique({
      where: { id: noteId, userId },
    });

    if (!note)
      throw new NotFoundException(`Note with provided id was not found`);

    const updatedNote = await this.prisma.note.update({
      where: { id: noteId },
      data: dto,
    });

    return updatedNote;
  }

  async deleteNote(userId: string, noteId: string) {
    const note = await this.prisma.note.findUnique({
      where: { id: noteId, userId },
    });

    if (!note)
      throw new NotFoundException(`Note with provided id was not found`);

    return this.prisma.note.delete({ where: { id: noteId } });
  }
}
