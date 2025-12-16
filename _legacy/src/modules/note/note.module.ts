import { Module } from '@nestjs/common';
import { NoteService } from './note.service';
import { NoteController } from './note.controller';
import { NoteTagController } from './note-tag/noteTag.controller';
import { NoteTagService } from './note-tag/noteTag.service';
import { TagModule } from '../tag/tag.module';

@Module({
  imports: [TagModule],
  controllers: [NoteController, NoteTagController],
  providers: [NoteService, NoteTagService],
})
export class NoteModule {}
