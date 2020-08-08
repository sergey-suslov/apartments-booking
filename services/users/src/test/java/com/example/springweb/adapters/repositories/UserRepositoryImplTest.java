package com.example.springweb.adapters.repositories;

import com.example.springweb.adapters.repositories.daos.UserDao;
import com.example.springweb.entities.UserEntity;
import org.junit.jupiter.api.Test;
import org.mockito.Mock;
import org.mockito.Mockito;

import static org.junit.jupiter.api.Assertions.*;

class UserRepositoryImplTest {

  @Test
  void saveUser() {

    String defaultEmail = "email@email.com";
    UserRepositoryImpl userRepository = new UserRepositoryImpl();

    UserDao mock = Mockito.mock(UserDao.class);
    userRepository.setUserDao(mock);

    // Setup input userEntity
    UserEntity userEntity = new UserEntity();
    userEntity.setEmail(defaultEmail);


    // Setup output userEntity
    UserEntity userEntityResult = new UserEntity();
    userEntityResult.setEmail(defaultEmail);
    userEntityResult.setId(1L);

    Mockito.when(mock.save(userEntity)).thenReturn(userEntityResult);

    UserEntity resultUserEntity = userRepository.saveUser(userEntity);
    assertEquals(resultUserEntity.getId(), 1L);
    assertEquals(resultUserEntity.getEmail(), defaultEmail);
  }
}