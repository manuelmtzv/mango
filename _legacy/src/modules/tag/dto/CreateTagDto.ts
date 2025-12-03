import {
  IsHexColor,
  IsNotEmpty,
  IsOptional,
  IsString,
  MaxLength,
  MinLength,
} from 'class-validator';
import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';

export class CreateTagDto {
  @ApiProperty({ example: 'Tag name' })
  @IsNotEmpty()
  @IsString()
  @MinLength(2)
  @MaxLength(20)
  name: string;

  @ApiPropertyOptional({ example: '#fff' })
  @IsOptional()
  @IsHexColor()
  color: string;
}
