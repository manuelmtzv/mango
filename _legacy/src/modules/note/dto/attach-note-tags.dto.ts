import { ApiProperty } from '@nestjs/swagger';
import { IsArray, IsNotEmpty, IsString } from 'class-validator';

export class AttachNoteTagsDto {
  @ApiProperty({
    description: 'Array of tag names to attach to the note',
    example: ['tag1', 'tag2'],
  })
  @IsArray()
  @IsString({ each: true })
  @IsNotEmpty({ each: true })
  tags: string[];
}
