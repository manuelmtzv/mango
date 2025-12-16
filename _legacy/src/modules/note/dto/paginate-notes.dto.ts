import { IsArray, IsNumber, IsOptional, IsString, Min } from 'class-validator';
import { ApiPropertyOptional } from '@nestjs/swagger';
import { Transform } from 'class-transformer';

export class PaginateNotesDto {
  @ApiPropertyOptional({ example: 10 })
  @IsOptional()
  @IsNumber()
  @Min(1)
  @Transform(({ value }) => parseInt(value))
  page?: number;

  @ApiPropertyOptional({ example: 5 })
  @IsOptional()
  @IsNumber()
  @Min(0)
  @Transform(({ value }) => parseInt(value))
  limit?: number;

  @ApiPropertyOptional({ example: 'Things to do' })
  @IsOptional()
  @IsString()
  search?: string;

  @ApiPropertyOptional({ example: ['tag1', 'tag2'] })
  @IsOptional()
  @IsArray()
  tags?: string[];
}
