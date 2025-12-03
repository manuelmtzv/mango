import { IsHexColor, IsOptional, IsString } from 'class-validator';
import { ApiPropertyOptional } from '@nestjs/swagger';

export class UpdateTagDto {
  @ApiPropertyOptional({ example: 'Tag name' })
  @IsOptional()
  @IsString()
  name: string;

  @ApiPropertyOptional({ example: '#000' })
  @IsOptional()
  @IsHexColor()
  color: string;
}
