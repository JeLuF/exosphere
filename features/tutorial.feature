@e2e
Feature: Following the tutorial

  As a person learning Exosphere
  I want that the whole tutorial works end to end
  So that I can follow along with the examples without getting stuck on bugs.

  AC:
  - all steps in the tutorial work when executed one after the other

  Notes:
  - The steps only do quick verifications.
    Full verifications are in the individual specs for the respective step.
  - You can not run individual scenarios here,
    you always have to run the whole feature.


  Scenario: tutorial
    ########################################
    # Printing the exosphere version
    ########################################
    When running "exo version" in my application directory
    Then it prints the current version in the terminal

    ########################################
    # Setting up the application
    ########################################
    Given I am in an empty folder
    When starting "exo init" in my application directory
    And entering into the wizard:
      | FIELD              | INPUT              |
      | AppName            | todo-app           |
      | ExocomVersion      | 0.27.0             |
    And waiting until the process ends
    Then my workspace contains the file "application.yml" with content:
      """
      name: todo-app

      local:
        dependencies:
          exocom:
            image: originate/exocom:0.27.0

      services:
      """
    And my workspace contains the empty directory ".exosphere/service_templates"
    And running "git init" in my application directory

    ########################################
    # Adding the html service
    ########################################
    Given I add the "exosphere-htmlserver-express" template
    When starting "exo add" in my application directory
    And entering into the wizard:
      | FIELD                         | INPUT                            |
      | template                      | 1                                |
      | serviceRole                   | html-server                      |
      | appName                       | test-app                         |
      | serviceType                   | public                           |
    And waiting until the process ends
    Then my application now contains the file "application.yml" with the content:
      """
      name: todo-app
      local:
        dependencies:
          exocom:
            image: originate/exocom:0.27.0
      services:
        html-server:
          location: ./html-server
      """
    And my application now contains the file "html-server/service.yml" with the content:
    """
    type: public

    development:
      scripts:
        run: node ./index.js
    """

    ########################################
    # adding the todo service
    ########################################
    Given I add the "exoservice-js-mongodb" template
    When starting "exo add" in my application directory
    And entering into the wizard:
      | FIELD                         | INPUT                    |
      | template                      | 1                        |
      | serviceRole                   | todo-service             |
      | serviceType                   | worker                   |
      | modelName                     | todo                     |
    And waiting until the process ends
    Then my application now contains the file "todo-service/service.yml" with the content:
      """
      type: worker

      dependency-data:
        exocom:
          receives:
            - todo.create
            - todo.create_many
            - todo.delete
            - todo.list
            - todo.read
            - todo.update
          sends:
            - todo.created
            - todo.created_many
            - todo.deleted
            - todo.listing
            - todo.details
            - todo.updated

      development:
        scripts:
          run: node src/server.js
          test: node_modules/cucumber/bin/cucumber.js

      local:
        dependencies:
          mongo:
            image: 'mongo:3.4.0'
            persist:
              - /data/db
      """

    ########################################
    # wiring up the html server to the todo service
    ########################################
    Given the file "html-server/app/controllers/index-controller.js":
      """
      class IndexController {

        constructor({send}) {
          this.send = send
        }

        index(req, res) {
          this.send('todo.list', {}, (messageName, payload) => {
            if (messageName === 'todo.listing') {
              res.render('index', {todos: payload})
            } else {
              res.sendStatus(500)
            }
          })
        }

      }

      module.exports = IndexController
      """
    And the file "html-server/app/views/index.pug":
      """
      extends layout

      block content

        h2 Exosphere Todos list
        p Your todos:
        ul
          each todo in todos
            li= todo.text

        h3 add a todo
        form(action="/todos" method="post")
          label text
          input(name="text")
          input(type="submit" value="add todo")
      """
    And the file "html-server/app/controllers/todos-controller.js":
      """
      class TodosController {

        constructor({send}) {
          this.send = send
        }

        create(req, res) {
          this.send('todo.create', req.body, () => {
            res.redirect('/')
          })
        }

      }
      module.exports = TodosController
      """
    And the file "html-server/app/routes.js":
      """
      module.exports = ({GET, resources}) => {
        GET('/', { to: 'index#index' })
        resources('todos', { only: ['create', 'destroy'] })
      }
      """
    And the file "html-server/service.yml":
      """
      type: public

      dependency-data:
        exocom:
          sends:
            - todo.create
            - todo.list
          receives:
            - todo.created
            - todo.listing

      development:
        port: 3000
        scripts:
          run: node ./index.js
      """
    When starting "exo run" in my application directory
    And waiting until I see "HTML server is running" in the terminal
    And waiting until I see "'todo-service' registered" in the terminal
    Then http://localhost:3000 displays:
      """
      Exosphere Todos list

      Your todos:
      """
    When adding a todo entry called "hello world" via the web application
    Then http://localhost:3000 displays:
      """
      Exosphere Todos list

      Your todos:

      * hello world
      """
    And I stop all running processes
