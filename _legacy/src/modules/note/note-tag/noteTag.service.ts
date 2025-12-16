import { Injectable, NotFoundException } from '@nestjs/common';
import { PrismaService } from '@modules/prisma/prisma.service';
import { TagService } from '@src/modules/tag/tag.service';
import { NoteService } from '@src/modules/note/note.service';

@Injectable()
export class NoteTagService {
  constructor(
    private readonly prisma: PrismaService,
    private readonly noteService: NoteService,
    private readonly tagService: TagService,
  ) {}

  async getNoteTags(userId: string, noteId: string) {
    return this.prisma.tag.findMany({
      where: { notes: { some: { id: noteId } } },
    });
  }

  async attachTag(userId: string, noteId: string, tagId: string) {
    const note = await this.prisma.note.findUnique({
      where: { id: noteId, userId },
    });

    if (!note)
      return new NotFoundException('Note with provided id was not found');

    return this.prisma.note.update({
      where: { id: noteId, userId },
      data: {
        tags: {
          connect: { id: tagId },
        },
      },
    });
  }

  async attachTags(userId: string, noteId: string, tagNames: string[]) {
    const tags = await this.tagService.findTagsOrCreate(userId, tagNames);

    const noteTags = (await this.noteService.getNote(userId, noteId)).tags;

    const disconectTags = noteTags.filter(
      (tag) => !tags.data.map((t) => t.id).includes(tag.id),
    );

    return this.prisma.note.update({
      where: { id: noteId, userId },
      data: {
        tags: {
          connect: tags.data.map((tag) => ({ id: tag.id })),
          disconnect: disconectTags.map((tag) => ({ id: tag.id })),
        },
      },
    });
  }

  async detachTag(userId: string, noteId: string, tagId: string) {
    const note = await this.prisma.note.findUnique({
      where: { id: noteId, userId },
    });

    if (!note)
      return new NotFoundException('Note with provided id was not found');

    return this.prisma.note.update({
      where: { id: noteId, userId },
      data: {
        tags: {
          disconnect: { id: tagId },
        },
      },
    });
  }
}
