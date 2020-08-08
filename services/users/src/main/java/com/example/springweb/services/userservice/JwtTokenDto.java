package com.example.springweb.services.userservice;

import lombok.Getter;
import lombok.Setter;

import java.util.Date;

@Getter
@Setter
public class JwtTokenDto {
  String jwt;
  Date expirationDate;
}
