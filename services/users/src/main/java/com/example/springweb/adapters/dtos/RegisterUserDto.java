package com.example.springweb.adapters.dtos;

import lombok.Getter;
import lombok.Setter;

import javax.persistence.Entity;
import javax.persistence.Id;
import javax.validation.constraints.Email;
import javax.validation.constraints.NotBlank;

@Entity
@Setter
@Getter
public class RegisterUserDto {
  @Email(message = "Email must be email")
  @Id
  private String email;

  @NotBlank
  private String password;
}
