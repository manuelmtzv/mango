import {
  Body,
  Controller,
  Delete,
  Get,
  HttpCode,
  HttpStatus,
  Param,
  Patch,
  Post,
  Query,
  UseGuards,
} from '@nestjs/common';
import { NoteService } from './note.service';
import { CreateNoteDto, PaginateNotesDto } from './dto';
import { IsObjectIdPipe } from '@src/pipes/mongoId/isObjectId.pipe';
import { GetUser } from '@modules/auth/decorator';
import { JwtGuard } from '@modules/auth/guard';
import { UpdateNoteDto } from './dto';
import { ApiBearerAuth, ApiTags } from '@nestjs/swagger';

@ApiTags('Notes')
@ApiBearerAuth('JWT-auth')
@UseGuards(JwtGuard)
@Controller('notes')
export class NoteController {
  constructor(private readonly noteService: NoteService) {}

  @Get()
  getNotes(
    @GetUser('id') userId: string,
    @Query() paginateDto: PaginateNotesDto,
  ) {
    return this.noteService.getNotes(userId, paginateDto);
  }

  @Get(':id')
  getNote(
    @GetUser('id') userId: string,
    @Param('id', IsObjectIdPipe) noteId: string,
  ) {
    return this.noteService.getNote(userId, noteId);
  }

  @Post()
  createNote(@GetUser('id') userId: string, @Body() dto: CreateNoteDto) {
    return this.noteService.createNote(userId, dto);
  }

  @Patch(':id')
  updateNote(
    @GetUser('id') userId: string,
    @Param('id', IsObjectIdPipe) noteId: string,
    @Body() dto: UpdateNoteDto,
  ) {
    return this.noteService.updateNote(userId, noteId, dto);
  }

  @HttpCode(HttpStatus.NO_CONTENT)
  @Delete(':id')
  deleteNote(
    @GetUser('id') userId: string,
    @Param('id', IsObjectIdPipe) noteId: string,
  ) {
    return this.noteService.deleteNote(userId, noteId);
  }
}
