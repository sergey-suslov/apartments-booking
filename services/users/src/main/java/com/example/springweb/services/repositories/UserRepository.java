package com.example.springweb.services.repositories;

import com.example.springweb.entities.UserEntity;

public interface UserRepository {
  UserEntity saveUser(UserEntity userEntity);
  UserEntity findUserByEmail(String email);
  UserEntity findUserById(long id);
}
