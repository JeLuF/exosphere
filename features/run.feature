Feature: running Exosphere applications

  As an Exosphere developer
  I want to have an easy way to run an application on my developer machine
  So that I can test my app locally.

  Rules:
  - run "exo run" in the directory of your application to run it
  - this command boots up all dependencies


  Scenario: running the "test" application
    When I start the "test" application
    Then my machine is running ExoComm
    And my machine is running the services:
      | NAME    |
      | web     |
      | users   |
