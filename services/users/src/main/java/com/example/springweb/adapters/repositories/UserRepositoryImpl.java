package com.example.springweb.adapters.repositories;

import com.example.springweb.adapters.repositories.daos.UserDao;
import com.example.springweb.entities.UserEntity;
import com.example.springweb.services.repositories.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.stereotype.Repository;

@Repository
public class UserRepositoryImpl implements UserRepository {
  UserDao userDao;

  @Autowired
  public void setUserDao(UserDao userDao) {
    this.userDao = userDao;
  }

  @Override
  public UserEntity saveUser(UserEntity userEntity) {
    return userDao.save(userEntity);
  }

  @Override
  public UserEntity findUserByEmail(String email) {
    return userDao.findByEmail(email);
  }

  @Override
  public UserEntity findUserById(long id) {
    return userDao.findById(id);
  }
}
