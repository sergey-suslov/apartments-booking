package com.example.springweb.adapters.repositories.daos;

import com.example.springweb.entities.UserEntity;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface UserDao extends CrudRepository<UserEntity, Long> {
  @Override
  <S extends UserEntity> S save(S entity);

  <S extends UserEntity> S findByEmail(String email);

  <S extends UserEntity> S findById(long id);
}
