# restapi-go-crud

Building a RESTFUL API with Go, PostgreSQL, GORM, Gin Framework

You can download this repository by using the green ``Clone or Download`` button on the right hand side of this page. This will present you with the option to either clone the repository using Git, or to download it as a zip file.

If you want to download it using git, copy paste the link that is presented to you, then run the following at your terminal:
 ```
git clone https://github.com/Rahul-Barick/restapi-go-crud.git
cd restapi-go-crud
```
# Prerequisites
- Go 1.21+
- Docker installed
- Makefile installed
  
# Steps to install & setup - run

1. COPY **.env.example** and create new file named as **.env**
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=admin
DB_NAME=internal_transfers
```
2. Start PostgreSQL Database
   - You already have PostgreSQL installed:-
     If Postgres is installed locally, create the database using your preferred tool (e.g: pgadmin or DBeaver)
   - You don't have PostgreSQL installed. Execute the command. **Note:- Make sure Docker is installed and running on your system.**
     This command spins up a Postgres container with the right credentials and database name.
     ```
     make docker-up
     ```
   - **Very IMP:- Run the Database queries of creation of table in **db/scripts** folder**
 3. Install Go Dependencies
    ```
    make deps
    ```
    Start the local server
    ```
    make run
    ```
    Note:- If you want help with more commands regarding this project
    ```
     make help
    ```
  4. If No MakeFile is installed, then install go dependencies
  ```
  go mod tidy
  ```
  Start the local server
  ```
  go run main.go
  ```
**Server will start running on PORT: 3000**
  Example of sample curl and **referenceId** is a idempotent key must he passed in headers of all POST API.
 **It must be passed ALWAYS UNIQUE ID as UUID**
  ```
  curl --location 'http://localhost:3000/accounts' \
--header 'Content-Type: application/json' \
--header 'referenceId: 4e176ce6-bd98-48c4-b11d-b7b0c0110934' \
--data '{
    "account_id": 123,
    "initial_balance": "100.23344"
}'
```

# Overview & Implementation
1. The project implements REST API's with proper validations, idempotency, and audit-compliant behavior using Go, Gin framework, and PostgreSQL. It is designed keeping clean architecture and best practices in mind.
2. Database Design:-
   - The system revolves around three main entities:
   - **Accounts**:- Represents a user’s account in the system that holds a monetary balance and is the foundation for all financial operations.
   - **Transactions**:- Represents a money transfer from one account to another. Ensures proper idempotency, validation, and consistency across the system.
   - **Ledger Entries**:-  Ledger is the Source of Truth. Stores dual-entry ledger rows per transaction - (CREDIT and DEBIT). It Ensures compliance-grade auditability of money movement.
3. This project is structured around Clean Architecture, encouraging separation of concerns.
   - **Handler Layer** (app/handler/) - Contains business entrypoints
   - **DTO Layer** (app/dto/) - Separates input/output schema & enforces strict validation using Gin binding and Go Validator tags
   - **Entities or Model** Layer (app/models/) - Entities like accounts, transaction, ledger_entries
   - **Utilities** (app/utils/) - Contains reusable logic
4. **Assumptions taken:-**
   - Accounts is an unique identifier in the core banking.
   - Every new account must have a positive initial balance, which will be credited as part of creation. Aligns with financial best practices to track all balance changes via transactions.
   - Ledger is the Source of Truth:- Any transaction must have a corresponding DEBIT and CREDIT entry. It Enables complete traceability, reversibility, and audit traceability.
   - Idempotency is Mandatory for All POST APIs. Prevents duplicate fund transfers during retries or network errors. Reduces backend complexity (no need to track duplicate attempts).
   - PostgreSQL is Used as a Strict Relational Store.  No NoSQL or in-memory DB is used — PostgreSQL handles all transactional and audit data.
   - Application Operates in UTC Only. It Ensures consistency across globally distributed systems.
   - No Soft Deletes:- Entities like accounts, transactions, and ledger_entries cannot be deleted. Financial records must be immutable for audit and compliance reasons. Hard deletes are also avoided to maintain foreign key integrity. If deletion is ever required, it must be handled via archival and not actual deletion.
5. **Technical implementation:-**
   - Idempotency Handling: Each transaction is uniquely identified via referenceId (UUID) to avoid duplicate processing.
   - Row-Level Locking: Uses gorm.Clauses(clause.Locking{Strength: "UPDATE"}) to handle concurrent balance changes safely. Transfers acquire row locks on source and destination accounts during balance modification & to prevent race conditions when multiple transactions happen on the same account.
   - Validation Layer: Strong validation via binding tags + custom validator integration.
   - Ledger-Based Double Entry: All money movement is traceable with DEBIT and CREDIT rows—ideal for audit logging.
   - Clean Code Structure: Handlers, DTOs, Models, Configs are logically separated. No tight coupling between DB or business logic.
   - All payload validation is performed at the request layer using DTOs.
   - Makefile is Used for Automation:- Common project tasks like installing dependecies, running the server and spinning up Docker are automated.
6. **Behaviour and purpose of ENTITIES**:-
   - ACCOUNTS:-
       - All incoming and outgoing transactions reference this with account number
       - Balance is updated only via transactions, never manually.
       - Always initialized with a non-zero initial_balance, recorded as a CREDIT ledger entry.
    - TRANSACTIONS:-
      - Transfer's money from one account to another. Ensures proper idempotency, validation, and consistency across the system.
      - Transaction is rejected if either account doesn’t exist or source lacks funds.
      - A successful transaction always generates: 1 debit entry (from source account) & 1 credit entry (to destination account)
    - LEDGER_ENTRIES:-
      - Every transaction creates two entries: DEBIT from source account & CREDIT to destination account
      - Maintains a financial audit trail.
      - Ledger entries is never updated or deleted. Immutable by design. 
   
