import { ApiProperty } from '@nestjs/swagger';
import { IsArray, IsString, MaxLength, MinLength } from 'class-validator';

export class CreateManyTagsDto {
  @ApiProperty({ example: ['Tag name 1', 'Tag name 2'] })
  @IsArray()
  @IsString({ each: true })
  @MinLength(2, { each: true })
  @MaxLength(20, { each: true })
  tags: string[];
}
