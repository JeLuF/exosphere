Feature: cleaning dangling Docker images

  As a developer using Exosphere
  I want to be able to easily clean up my Docker workspace
  So that unused images and volumes do not take up space on my Docker VM

  Rules:
  - run "exo clean" in the terminal in any directory to remove dangling
    Docker images and volumes on your machine
  - this command does not remove non-dangling Docker images/volumes


  Background:
    Given I am in the root directory of the "clean-containers" example application


  Scenario: cleaning a machine with both dangling and non-dangling Doker images
    Given my machine has both dangling and non-dangling Docker images and volumes
    When running "exo clean" in my application directory
    Then it prints "Removing dangling images" in the terminal
    And it prints "Removing dangling volumes" in the terminal
    And it has non-dangling images
    And it does not have dangling images
    And it does not have dangling volumes

  Scenario: cleaning a machine with both dangling and non-dangling Doker images from a service directory
    Given my machine has both dangling and non-dangling Docker images and volumes
    When running "exo clean" in the "test-service" directory
    Then it prints "Removing dangling images" in the terminal
    And it prints "Removing dangling volumes" in the terminal
    And it has non-dangling images
    And it does not have dangling images
    And it does not have dangling volumes


  Scenario: cleaning a machine with running application containers
    Given my machine has running application containers
    And my machine has running third party containers
    When running "exo clean" in my application directory
    Then it prints "Removing application and test containers" in the terminal
    And it prints "Stopping application-service" in the terminal
    And it prints "Removing application-service" in the terminal
    And it removes application and test containers
    And it does not stop any third party containers


  Scenario: cleaning a machine with running test containers
    Given my machine has running test containers
    And my machine has running third party containers
    When running "exo clean" in my application directory
    Then it prints "Removing application and test containers" in the terminal
    And it prints "Stopping application-service" in the terminal
    And it prints "Removing application-service" in the terminal
    And it removes application and test containers
    And it does not stop any third party containers


  Scenario: cleaning a machine with stopped application containers
    Given my machine has stopped application containers
    And my machine has running third party containers
    When running "exo clean" in my application directory
    Then it prints "Removing application and test containers" in the terminal
    Then it prints "Removing application-service" in the terminal
    And it removes application and test containers
    And it does not stop any third party containers


  Scenario: cleaning a machine with stopped test containers
    Given my machine has stopped test containers
    And my machine has running third party containers
    When running "exo clean" in my application directory
    Then it prints "Removing application and test containers" in the terminal
    Then it prints "Removing application-service" in the terminal
    And it removes application and test containers
    And it does not stop any third party containers
