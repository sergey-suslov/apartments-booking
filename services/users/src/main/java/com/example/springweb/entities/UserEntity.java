package com.example.springweb.entities;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import lombok.Getter;
import lombok.Setter;

import javax.persistence.*;
import javax.validation.constraints.Email;
import javax.validation.constraints.NotBlank;

@Entity
@Table(name = "users")
@Setter
@Getter
@JsonIgnoreProperties(value = {"hash", "salt"})
public class UserEntity {

  @Id
  @GeneratedValue(strategy = GenerationType.SEQUENCE, generator = "user_seq")
  private Long id;

  @Email(message = "Email must be email")
  @NotBlank
  private String email;

  private byte[] hash;
  private byte[] salt;
}
