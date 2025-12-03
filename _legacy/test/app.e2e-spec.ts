import { Test } from '@nestjs/testing';
import { INestApplication, ValidationPipe } from '@nestjs/common';
import * as request from 'supertest';
import { AppModule } from '@src/app.module';
import { PrismaService } from '@src/modules/prisma/prisma.service';
import { RegisterDto } from '@src/modules/auth/dto/register.dto';

describe('App (e2e)', () => {
  let app: INestApplication;
  let prisma: PrismaService;

  beforeAll(async () => {
    const testModule = await Test.createTestingModule({
      imports: [AppModule],
    }).compile();

    app = testModule.createNestApplication();

    app.useGlobalPipes(
      new ValidationPipe({
        whitelist: true,
      }),
    );

    await app.init();
    await app.listen(4000);

    prisma = app.get(PrismaService);
    await prisma.cleanDb();
  });

  describe('Auth', () => {
    const registerDto: RegisterDto = {
      email: 'manuel.mtzv816@gmail.com',
      password: 'passWord64$',
      username: 'manuel.mtzv',
      name: 'Manuel',
    };

    describe('Register', () => {
      it('Should throw an error if email is empty', () => {
        return request(app.getHttpServer())
          .post('/auth/register')
          .send({
            ...registerDto,
            email: '',
          })
          .expect(400);
      });

      it('Should throw an error password is empty', () => {
        return request(app.getHttpServer())
          .post('/auth/register')
          .send({
            ...registerDto,
            password: '',
          })
          .expect(400);
      });

      it('Should throw an error if username is empty', () => {
        return request(app.getHttpServer())
          .post('/auth/register')
          .send({
            ...registerDto,
            username: '',
          })
          .expect(400);
      });

      it('Should successfully register a user', () => {
        return request(app.getHttpServer())
          .post('/auth/register')
          .send(registerDto)
          .expect(201);
      });
    });

    describe('Login', () => {
      it('Should throw an error if identifier or/and password are empty', () => {
        return request(app.getHttpServer())
          .post('/auth/login')
          .send({
            identifier: '',
            password: '',
          })
          .expect(400);
      });

      it('Should throw an error if password is not correct', () => {
        return request(app.getHttpServer())
          .post('/auth/login')
          .send({
            identifier: registerDto.email,
            password: 'password',
          })
          .expect(403);
      });

      it("Should throw an error if user doesn't exist", () => {
        return request(app.getHttpServer())
          .post('/auth/login')
          .send({
            identifier: 'johndoe123',
            password: registerDto.password,
          })
          .expect(403);
      });

      it('Should successfully login a user', () => {
        return request(app.getHttpServer())
          .post('/auth/login')
          .send({
            identifier: registerDto.email,
            password: registerDto.password,
          })
          .expect(200);
      });
    });
  });

  afterAll(async () => {
    await app.close();
  });
});
