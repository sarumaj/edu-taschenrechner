Feature: Example test feature to demonstrate BDT

    Scenario: I enter 1.3+(12*-7)+1
        When I press following buttons: "1 . 3 + () 1 2 × - - - 7 () + 1"
        Then I get following result: "1.3+(12×-7)+1_"

    Scenario: I calculate 1.3+(12*-7)+1
        When I press following buttons: "1 . 3 + () 1 2 × - 7 () + 1 ="
        Then I get following result: "-81.7"

    Scenario: I enter 1.3+(12*-7)+1 and then use memory cell to enter ANS*6/7
        When I press following buttons: "1 . 3 + () 1 2 × - 7 () + 1 = 6 ÷ 7"
        Then I get following result: "ANS×6÷7_"

    Scenario: I enter 1.3+(12*-7)+1 and then use memory cell to calculate ANS*6/7
        When I press following buttons: "1 . 3 + () 1 2 × - 7 () + 1 = 6 ÷ 7 ="
        Then I get following result: "-70.02857142857144"