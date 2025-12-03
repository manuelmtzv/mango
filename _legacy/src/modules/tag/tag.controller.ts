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
  UseGuards,
} from '@nestjs/common';
import { TagService } from './tag.service';
import { CreateTagDto } from './dto/CreateTagDto';
import { UpdateTagDto } from './dto/UpdateTag';
import { GetUser } from '@modules/auth/decorator';
import { JwtGuard } from '@modules/auth/guard';
import { IsObjectIdPipe } from '@src/pipes/mongoId/isObjectId.pipe';
import { ApiBearerAuth, ApiTags } from '@nestjs/swagger';
import { CreateManyTagsDto } from './dto/CreateManyTags.dto';

@ApiTags('Tags')
@ApiBearerAuth('JWT-auth')
@UseGuards(JwtGuard)
@Controller('tags')
export class TagController {
  constructor(private readonly tagService: TagService) {}

  @Get()
  getTags(@GetUser('id') userId: string) {
    return this.tagService.getTags(userId);
  }

  @Get(':id')
  getTag(
    @GetUser('id') userId: string,
    @Param('id', IsObjectIdPipe) tagId: string,
  ) {
    return this.tagService.getTag(userId, tagId);
  }

  @Post()
  createTag(@GetUser('id') userId: string, @Body() dto: CreateTagDto) {
    return this.tagService.createTag(userId, dto);
  }

  @Post('find-or-create')
  findTagsOrCreate(
    @GetUser('id') userId: string,
    @Body() createManyTagsDto: CreateManyTagsDto,
  ) {
    return this.tagService.findTagsOrCreate(userId, createManyTagsDto.tags);
  }

  @Patch(':id')
  updateTag(
    @GetUser('id') userId: string,
    @Param('id', IsObjectIdPipe) tagId: string,
    @Body() dto: UpdateTagDto,
  ) {
    return this.tagService.updateTag(userId, tagId, dto);
  }

  @HttpCode(HttpStatus.NO_CONTENT)
  @Delete(':id')
  deleteTag(
    @GetUser('id') userId: string,
    @Param('id', IsObjectIdPipe) tagId: string,
  ) {
    return this.tagService.deleteTag(userId, tagId);
  }
}
