import { Injectable, PipeTransform, BadRequestException } from '@nestjs/common';
import { isMongoId } from 'class-validator';

@Injectable()
export class IsObjectIdPipe implements PipeTransform {
  transform(value: string) {
    if (!isMongoId(value)) {
      throw new BadRequestException(`Invalid ObjectId`);
    }
    return value;
  }
}
