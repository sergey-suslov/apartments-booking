package com.example.springweb.configurations;

import com.example.springweb.adapters.userservice.UserServiceGRPCController;
import io.grpc.ServerBuilder;
import lombok.SneakyThrows;
import org.apache.tomcat.dbcp.dbcp2.BasicDataSource;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.Configuration;
import org.springframework.orm.jpa.JpaTransactionManager;
import org.springframework.orm.jpa.JpaVendorAdapter;
import org.springframework.orm.jpa.LocalContainerEntityManagerFactoryBean;
import org.springframework.orm.jpa.vendor.HibernateJpaVendorAdapter;
import org.springframework.transaction.PlatformTransactionManager;
import org.springframework.transaction.annotation.EnableTransactionManagement;

import javax.annotation.PostConstruct;
import javax.sql.DataSource;
import java.util.Properties;

@Configuration
@ComponentScan(value = "com.example.springweb.adapters.userservice")
public class UserServiceGrpcControllerConfiguration {

  public UserServiceGRPCController userServiceGRPCController;

  @Autowired
  public void setUserServiceGRPCController(UserServiceGRPCController userServiceGRPCController) {
    this.userServiceGRPCController = userServiceGRPCController;
  }

  @PostConstruct
  public void getUserServiceGrpcController() {
    this.startUserServiceGrpc();
  }

  @SneakyThrows
  private void startUserServiceGrpc() {
    ServerBuilder.forPort(50051).addService(this.userServiceGRPCController).build().start();
  }
}
