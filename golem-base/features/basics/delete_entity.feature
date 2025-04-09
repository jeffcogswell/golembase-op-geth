Feature: deleting entities

  Scenario: deleting an existing entity
    Given I have created an entity
    When I submit a transaction to delete the entity
    Then the entity should be deleted
    And the number of entities should be 0
    And the list of all entities should be empty

  Scenario: deleting entity after it has been updated
    Given I have created an entity
    When I submit a transaction to update the entity, changing the annotations
    And I submit a transaction to delete the entity
    Then the entity should be deleted
    And the number of entities should be 0
    And the list of all entities should be empty
    And the owner should not have any entities
