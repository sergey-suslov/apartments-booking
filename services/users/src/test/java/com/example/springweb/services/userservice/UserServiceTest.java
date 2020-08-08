package com.example.springweb.services.userservice;

import com.example.springweb.entities.UserEntity;
import com.example.springweb.services.repositories.UserRepository;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;
import org.springframework.http.HttpStatus;
import org.springframework.web.server.ResponseStatusException;

import java.util.Date;

import static org.junit.jupiter.api.Assertions.*;

class UserServiceTest {
  @Test
  void registerUser() {
    String email = "email@mail.com";
    byte[] encodedPassword = new byte[]{1, 2, 3};

    UserService userService = new UserService();

    // Мок SHA256PasswordEncoder, возвращает encodedPassword
    SHA256PasswordEncoder sha256PasswordEncoder = Mockito.mock(SHA256PasswordEncoder.class);
    Mockito.when(sha256PasswordEncoder.encode(Mockito.anyString(), Mockito.any(byte[].class))).thenReturn(encodedPassword);
    userService.setPasswordEncoder(sha256PasswordEncoder);

    // Конфигурация объекта ответа на userRepository.save
    long defaultId = 1L;
    UserEntity userEntityWithId = new UserEntity();
    userEntityWithId.setEmail(email);
    userEntityWithId.setHash(encodedPassword);
    userEntityWithId.setId(defaultId);

    // Мок userRepository
    UserRepository userRepository = Mockito.mock(UserRepository.class);
    Mockito.doAnswer(invocation -> {
      // Получаем userEntity как первый аргумент при вызове userRepository.save
      UserEntity userEntity = invocation.getArgument(0);

      // Проверка экземпляра отданного на сохранение в БД
      // Email должен совпадать с отданным на вход сервису
      assertEquals(userEntity.getEmail(), email);
      // Пароль должен быть равен тому, который вернет мок SHA256PasswordEncoder
      assertEquals(userEntity.getHash(), encodedPassword);
      return userEntityWithId;
    }).when(userRepository).saveUser(Mockito.any(UserEntity.class));
    userService.setUserRepository(userRepository);

    UserEntity registerUser = userService.registerUser(email, "password");
    assertEquals(registerUser, userEntityWithId);
  }

  @Test
  void auth() {
    String defaultPassword = "password";
    String defaultEmail = "email@mail.com";
    String defaultJwt = "jwt-token-mock";
    byte[] encodedPassword = new byte[]{1, 2, 3};

    long defaultId = 1L;
    UserEntity userEntityWithId = new UserEntity();
    userEntityWithId.setEmail(defaultEmail);
    userEntityWithId.setHash(encodedPassword);
    userEntityWithId.setId(defaultId);

    UserService userService = new UserService();

    UserRepository userRepository = Mockito.mock(UserRepository.class);
    Mockito.when(userRepository.findUserByEmail(defaultEmail)).thenReturn(userEntityWithId);

    SHA256PasswordEncoder passwordEncoder = Mockito.mock(SHA256PasswordEncoder.class);
    Mockito.when(passwordEncoder.match(Mockito.any(), Mockito.any(), Mockito.any())).thenReturn(true);

    JwtEncoder jwtEncoder = Mockito.mock(JwtEncoder.class);
    JwtTokenDto defaultJwtTokenDto = new JwtTokenDto();
    defaultJwtTokenDto.setJwt(defaultJwt);
    defaultJwtTokenDto.setExpirationDate(new Date());
    Mockito.when(jwtEncoder.generateToken(userEntityWithId)).thenReturn(defaultJwtTokenDto);

    userService.setUserRepository(userRepository);
    userService.setPasswordEncoder(passwordEncoder);
    userService.setJwtEncoder(jwtEncoder);

    JwtTokenDto jwtTokenDto = userService.auth(defaultEmail, defaultPassword);
    assertEquals(jwtTokenDto, defaultJwtTokenDto);
  }

  @Test
  void authException401() {
    String defaultPassword = "password";
    String defaultEmail = "email@mail.com";
    String defaultJwt = "jwt-token-mock";
    byte[] encodedPassword = new byte[]{1, 2, 3};

    long defaultId = 1L;
    UserEntity userEntityWithId = new UserEntity();
    userEntityWithId.setEmail(defaultEmail);
    userEntityWithId.setHash(encodedPassword);
    userEntityWithId.setId(defaultId);

    UserService userService = new UserService();

    UserRepository userRepository = Mockito.mock(UserRepository.class);
    Mockito.when(userRepository.findUserByEmail(defaultEmail)).thenReturn(userEntityWithId);

    SHA256PasswordEncoder passwordEncoder = Mockito.mock(SHA256PasswordEncoder.class);
    Mockito.when(passwordEncoder.match(Mockito.any(), Mockito.any(), Mockito.any())).thenReturn(false);

    JwtEncoder jwtEncoder = Mockito.mock(JwtEncoder.class);
    JwtTokenDto defaultJwtTokenDto = new JwtTokenDto();
    defaultJwtTokenDto.setJwt(defaultJwt);
    defaultJwtTokenDto.setExpirationDate(new Date());
    Mockito.when(jwtEncoder.generateToken(userEntityWithId)).thenReturn(defaultJwtTokenDto);

    userService.setUserRepository(userRepository);
    userService.setPasswordEncoder(passwordEncoder);
    userService.setJwtEncoder(jwtEncoder);

    ResponseStatusException responseStatusException = assertThrows(ResponseStatusException.class, () -> {
      JwtTokenDto jwtTokenDto = userService.auth(defaultEmail, defaultPassword);
    });

    String responseStatusExceptionMessage = responseStatusException.getMessage();
    assertNotNull(responseStatusExceptionMessage);
    assertTrue(responseStatusExceptionMessage.contains("Wrong email or password"));
    assertEquals(responseStatusException.getStatus(), HttpStatus.UNAUTHORIZED);
  }
}