> INTERNAL APPROVAL FORM

**Project Name:** Agent Incentive, Commission and Producer Management

**Version: 1.0**

**Submitted on:**

  ------------------------------------------------------------------------
               **Name**                                 **Date**
  ------------ ---------------------------------------- ------------------
  **Approved                                            
  By:**                                                 

  **Reviewed                                            
  By:**                                                 

  **Prepared                                            
  By: **                                                
  ------------------------------------------------------------------------

> VERSION CONTROL LOG

  -------------------------------------------------------------------------------
  **Version**   **Date**   **Prepared     **Remarks**
                           By**           
  ------------- ---------- -------------- ---------------------------------------
  **1**                                   

                                          

                                          

                                          

                                          
  -------------------------------------------------------------------------------

# Table of Contents {#table-of-contents .TOC-Heading}

[**1. Executive Summary** [4](#executive-summary)](#executive-summary)

[**2. Project Scope** [4](#project-scope)](#project-scope)

[**3. Business Requirements**
[4](#business-requirements)](#business-requirements)

[**4. Functional Requirements Specification**
[6](#functional-requirements-specification)](#functional-requirements-specification)

[4.1 Agent Onboarding [6](#agent-onboarding)](#agent-onboarding)

[4.1.1 New Profile Options
[6](#new-profile-options)](#new-profile-options)

[4.1.2 Enter Profile Details Page
[6](#enter-profile-details-page)](#enter-profile-details-page)

[4.1.3 Select New Advisor Coordinator Page
[8](#select-new-advisor-coordinator-page)](#select-new-advisor-coordinator-page)

[4.2 Agent Profile Management
[8](#agent-profile-management)](#agent-profile-management)

[4.2.1 Agent Search [8](#agent-search)](#agent-search)

[4.2.2 Agent Profile Maintenance Page
[9](#agent-profile-maintenance-page)](#agent-profile-maintenance-page)

[4.2.3 License Management Page
[10](#license-management-page)](#license-management-page)

[4.2.3 Agent Termination Page
[10](#agent-termination-page)](#agent-termination-page)

[4.3 Agent Commission Management
[11](#agent-commission-management)](#agent-commission-management)

[4.3.1 Commission Rate Table View Page
[11](#commission-rate-table-view-page)](#commission-rate-table-view-page)

[4.3.2 Commission History Search Page
[11](#commission-history-search-page)](#commission-history-search-page)

[4.3.3 Commission Processing Steps:
[11](#commission-processing-steps)](#commission-processing-steps)

[**5. Appendices** [15](#appendices)](#appendices)

## **1. Executive Summary**

To provide the functional and non-functional requirements for developing
an Agent Onboarding & Commission Management System similar to PMACS for
India Post PLI. The system will support onboarding and administration of
Agents, Advisor Coordinators, Field Officers, and Departmental
Employees; manage licensing and status lifecycle; and calculate,
approve, and disburse commissions.

## **2. Project Scope**

This system will support:

- Agent onboarding (Advisor, Coordinator, Field Officer, Departmental
  Employee)

- Agent profile management

- Commission rate setup and processing

- Trail and Final Incentive statement generation and disbursement

- Licensing and termination workflows

## **3. Business Requirements**

+-------------+---------------+-----------------------------------------------+
| Requirement | Functionality | Requirement                                   |
| ID          |               |                                               |
+:============+===============+:==============================================+
| FS_IC_001   | Agent         | The system shall allow creation of a new      |
|             | Onboarding    | Advisor profile, which must be linked to an   |
|             |               | existing Advisor Coordinator already present  |
|             |               | in the system.                                |
+-------------+               +-----------------------------------------------+
| FS_IC_002   |               | The system shall support creation of a new    |
|             |               | Advisor Coordinator profile, including        |
|             |               | assignment to a specific circle and division. |
+-------------+               +-----------------------------------------------+
| FS_IC_003   |               | The system shall enable onboarding of         |
|             |               | Departmental Employees by auto-populating     |
|             |               | profile data using a valid Employee ID from   |
|             |               | the HRMS system.                              |
+-------------+               +-----------------------------------------------+
| FS_IC_004   |               | The system shall allow onboarding of Field    |
|             |               | Officers either by auto-fetching data using   |
|             |               | Employee ID or through manual entry.          |
+-------------+---------------+-----------------------------------------------+
| FS_IC_005   | Agent Profile | The system shall provide a search interface   |
|             | Management    | to locate agent profiles using Agent ID,      |
|             |               | Name, PAN, or Mobile Number.                  |
+-------------+               +-----------------------------------------------+
| FS_IC_006   |               | The system shall display a dashboard view of  |
|             |               | the agent profile with editable sections for  |
|             |               | each data category.                           |
+-------------+               +-----------------------------------------------+
| FS_IC_007   |               | The system shall allow authorized users to    |
|             |               | update the agent's name with proper           |
|             |               | validation and audit logging.                 |
+-------------+               +-----------------------------------------------+
| FS_IC_008   |               | The system shall allow updating of PAN        |
|             |               | information with format validation and        |
|             |               | uniqueness checks.                            |
+-------------+               +-----------------------------------------------+
| FS_IC_009   |               | The system shall support status updates for   |
|             |               | agents (e.g., Active, Suspended, Terminated)  |
|             |               | along with mandatory reason entry.            |
+-------------+               +-----------------------------------------------+
| FS_IC_010   |               | The system shall allow updating of personal   |
|             |               | information including date of birth, gender,  |
|             |               | and marital status.                           |
+-------------+               +-----------------------------------------------+
| FS_IC_011   |               | The system shall support addition and         |
|             |               | modification of distribution channel details  |
|             |               | with effective dates.                         |
+-------------+               +-----------------------------------------------+
| FS_IC_012   |               | The system shall allow entry and update of    |
|             |               | external identification numbers and their     |
|             |               | sources.                                      |
+-------------+               +-----------------------------------------------+
| FS_IC_013   |               | The system shall support assignment and       |
|             |               | modification of product class information     |
|             |               | linked to the agent.                          |
+-------------+               +-----------------------------------------------+
| FS_IC_014   |               | The system shall allow entry and update of    |
|             |               | multiple address types: Official, Permanent,  |
|             |               | and Communication.                            |
+-------------+               +-----------------------------------------------+
| FS_IC_015   |               | The system shall support entry and update of  |
|             |               | phone numbers including official/resident     |
|             |               | landline and mobile numbers.                  |
+-------------+               +-----------------------------------------------+
| FS_IC_016   |               | The system shall allow entry and update of    |
|             |               | email addresses categorized as official,      |
|             |               | permanent, and communication.                 |
+-------------+               +-----------------------------------------------+
| FS_IC_017   |               | The system shall support assignment of        |
|             |               | authority types and validity periods for      |
|             |               | agents.                                       |
+-------------+               +-----------------------------------------------+
| FS_IC_018   |               | The system shall allow entry and update of    |
|             |               | insurance licensing details and generate      |
|             |               | automated reminders at 1 month, 15 days, 7    |
|             |               | days, and on the day of expiry.               |
+-------------+               +-----------------------------------------------+
| FS_IC_019   |               | The system shall enforce license renewal      |
|             |               | rules: first renewal after 1 year, subsequent |
|             |               | renewals every 3 years.                       |
+-------------+               +-----------------------------------------------+
| FS_IC_020   |               | The system shall automatically deactivate an  |
|             |               | advisor code if the license renewal date has  |
|             |               | elapsed, or allow manual cancellation via the |
|             |               | License Update Entry page.                    |
+-------------+---------------+-----------------------------------------------+
| FS_IC_021   | Advisor       | The system shall allow termination of an      |
|             | Termination   | advisor profile with mandatory entry of       |
|             |               | termination reason and effective date.        |
+-------------+---------------+-----------------------------------------------+
| FS_IC_022   | Commission    | The system shall provide a Commission Rate    |
|             | Processing    | Table setup interface with fields for Rate,   |
|             |               | Policy Duration (Months), Product Type,       |
|             |               | Product Plan Code, Agent Type, and Policy     |
|             |               | Term (Years).                                 |
+-------------+               +-----------------------------------------------+
| FS_IC_023   |               | The system shall allow searching of           |
|             |               | commission history by policy number and agent |
|             |               | ID.                                           |
+-------------+               +-----------------------------------------------+
| FS_IC_024   |               | The system shall support execution of monthly |
|             |               | Commission Calculation Batch jobs to compute  |
|             |               | agent commissions.                            |
+-------------+               +-----------------------------------------------+
| FS_IC_025   |               | The system shall support automatic generation |
|             |               | of Trial Statements via batch job based on    |
|             |               | policies sold.                                |
+-------------+               +-----------------------------------------------+
| FS_IC_026   |               | The system shall provide a page to view       |
|             |               | generated Trial Statements with agent-wise    |
|             |               | commission details.                           |
+-------------+               +-----------------------------------------------+
| FS_IC_027   |               | The system shall provide a Manual Trial       |
|             |               | Statement Generation page with fields for     |
|             |               | Processing Unit, Statement Format, Max        |
|             |               | Statement Due Date, Max Transaction Effective |
|             |               | Date, Max Process Date, Statement Date,       |
|             |               | Contract Holder, Advisor Coordinator, and     |
|             |               | Carrier.                                      |
+-------------+               +-----------------------------------------------+
| FS_IC_028   |               | The system shall provide an Approving Trial   |
|             |               | Statement page that displays commission       |
|             |               | amounts and allows full or partial            |
|             |               | disbursement approval.                        |
+-------------+               +-----------------------------------------------+
| FS_IC_029   |               | The system shall support execution of Final   |
|             |               | Incentive Statement Generation batch job      |
|             |               | after trial statement approval.               |
+-------------+               +-----------------------------------------------+
| FS_IC_030   |               | The system shall provide a Final Statements   |
|             |               | page displaying final commission amounts for  |
|             |               | agents.                                       |
+-------------+               +-----------------------------------------------+
| FS_IC_031   |               | The system shall provide a Disbursement       |
|             |               | Details page to input cheque or EFT           |
|             |               | information.                                  |
+-------------+               +-----------------------------------------------+
| FS_IC_032   |               | The system shall support automatic            |
|             |               | disbursement of commission amounts based on   |
|             |               | final statements, with immediate processing   |
|             |               | for cheque and queued processing for EFT.     |
+-------------+---------------+-----------------------------------------------+

## **4. Functional Requirements Specification**

## 4.1 Agent Onboarding

### 4.1.1 New Profile Options

- **Purpose:** To select the type of Agent that needs to be onboarded.

- **Fields & Rules:**

  - Agent Type: dropdown: Options are Advisor, Advisor Coordinator,
    Departmental Employee, Field Officer

  - Employee Number: Textbox: Mandatory for Departmental Employee,
    optional for Field Officer, Not Applicable for others.

  - Person Type: Dropdown: Options are Individual, Corporate/Group.

  - Advisor Undergoing Training: Checkbox: Default Unchecked

  - Continue: button

### 4.1.2 Enter Profile Details Page

- **Purpose:** This page will get displayed after selecting new profile
  options and clicking continue button.

- **Fields & Rules:**

  - Profile Type: Dropdown

  - Office Type: Dropdown

  - Office Code: Textbox

  - Advisor Sub-Type: Dropdown

  - Effective Date: Calendar

  - Distribution Channel: Multiselect Dropdown: Options are India Post.

  - Product Class: Multiselect Dropdown: Options are PLI, RPLI

  - Title: Dropdown

  - First Name: Textbox

  - Middle Name: Textbox

  - Last Name: Textbox

  - Gender: Dropdown: Options are Male, Female, Other

  - Date of Birth: Calendar

  - Category: Dropdown

  - Marital Status: Dropdown

  - Aadhar Number: Textbox

  - PAN: Textbox

  - Designation/Rank: Dropdown

  - Service Number: Textbox

  - Professional Title: Dropdown

  - Address:

    - Address Type: Dropdown: Options as Official, Permanent,
      Communication.

    - Address Line1-

    - Address Line2-

    - Village-

    - Taluka-

    - City-

    - District-

    - State-

    - Country-

    - Pin Code-

  - Phone:

  - Email:

  - Bank Account#:

  - Bank IFSC Code:

  - Superior Advisor: This will open Select New Advisor Coordinator
    Page, and the user will return to this page after selection of
    advisors.

  - Office Affiliation: Textbox: Input Affiliated Office Code

  ------------------------------------------------------------------------
  Serial     Condition           Error Message           Required Action
  Number                                                 
  ---------- ------------------- ----------------------- -----------------
  1          If the **Profile    Please select a Profile Users must select
             type** is not       Type.                   the Profile Type
             selected and                                from the
             **Continue** button                         drop-down list.
             is pressed.                                 

  2          If the PAN number   PAN number entered      Users must enter
             entered already     already exists for      the PAN number
             exists for some     another advisor's       which does not
             other profile.      profile and cannot be   exist for another
                                 for this profile.       advisor's
                                                         profile.

  3          PAN number should   Please enter a 10 digit Users must enter
             be of 10            Permanent Account       the 10 digit PAN
             characters. If the  Number (PAN).           number.
             PAN number length                           
             doesn't match.                              

  4          PAN number should   Please enter correct    Users must enter
             be entered in the   PAN.                    the PAN number as
             standard format as                          per the standard
             shown above. If the                         format.
             PAN doesn't match                           
             with the format as                          
             already defined in                          
             the system.                                 

  5          If the Last name is Please enter a Last     Users must enter
             not entered and     name.                   the last name of
             **Continue** button                         the Advisor.
             is pressed.                                 

  6          If the First name   Please enter a First    Users must enter
             is not entered and  name.                   the first name of
             **Continue** button                         the Advisor.
             is pressed.                                 

  7          If the **Date of    Please enter a valid    Users must enter
             Birth** is not      Date of Birth           the valid Date Of
             entered and                                 Birth.
             **Continue** button                         
             is pressed.                                 
  ------------------------------------------------------------------------

### 4.1.3 Select New Advisor Coordinator Page

- **Purpose:** If the new profile reports to an Advisor Coordinator
  already present in system, then select Agent type as 'Advisor' and
  click the Continue button to move the user to the Select New Advisor
  Coordinator page.

- **Fields & Rules:**

  - AC Profile#: text box

  - AC Profile Name: text box

  - Search: button

- Clicking the search button should display the following table with
  auto-populated agent details:

  - AC Profile#

  - AC Profile Name

  - Profile Type

  - Status

  - Person: Individual or Corporate/Group

  - Action: Link for selecting the Advisor Coordinator

  -----------------------------------------------------------------------
  Serial Number Error Message                     Required Action
  ------------- --------------------------------- -----------------------
  1             Your selected criteria did not    Users must change the
                return any rows. Please change    selection
                your selections and try again.    

  -----------------------------------------------------------------------

## 4.2 Agent Profile Management

### 4.2.1 Agent Search

- **Purpose:** To search an Existing Agent for performing some action
  like View / Edit / Terminate / Suspend / viewing Commission History.

- **Fields:**

  - Agent ID: Textbox

  - Last Name: Textbox

  - First Name: Textbox

  - PAN: Textbox

  - Mobile Number: Textbox

  - Status: Textbox

  - Superior Advisor Code: Textbox

  - Office ID: Textbox

  - Advisor Undergoing Training: Checkbox

  - Search: button

- **Business Rules**:

  - Results displayed in a table with clickable rows.

  - Clicking on any Agent in table will open that Agent Profile
    Maintenance Page for that agent.

  - 'Export to Excel' option should be present for the table details to
    be exported in excel format.

### 4.2.2 Agent Profile Maintenance Page

- **Purpose:** Displays agent details with option to edit each section.

- **Fields & Rules:** It should display all the information about the
  agent with 'Update' link in each section to update the details of that
  section. The sections that need to be displayed are:

  - Advisor Name Section

  - PAN Number Information

  - Status Information: Add additional 'Change Status To' dropdown for
    updating the Agent Status to Expired/Suspended/ Terminated.

  - Personal Information

  - Distribution Channel

  - External Number

  - Product Class

![A screenshot of a computer AI-generated content may be
incorrect.](media/image1.png){width="6.268055555555556in"
height="3.451388888888889in"}

### 4.2.3 License Management Page

- **Purpose:** This page will display the list of licenses the current
  agent have along with the option to Add or Delete License.

- **Fields:**

  - License Line: Dropdown: Option as Life

  - License Type: Dropdown

  - License Number: Textbox

  - Resident Status: Dropdown: Option as Resident, Non-Resident

  - License Date: Calendar

  - Renewal Date: Calendar

  - Authority Date: Calendar

  - Submit: Button

  - Update Renewal Date: Button

  - Delete License: Button

- **Business Rules:**

  - First License Renewal notice will be generated for the license
    renewal 1 month before the License expiry date.

  - Second License Renewal notice will be generated for the license
    renewal 15 days before the License expiry date, if not renewed
    already.

  - Third License Renewal notice will be generated for the license
    renewal 7 days before the License expiry date, if not renewed
    already.

  - Final License Renewal notice will be generated for the license
    renewal on the License expiry date, if not renewed already.

### 4.2.3 Agent Termination Page

- Search an Agent.

- Move to 'Agent Profile Maintenance Page' for the Agent.

- For terminating a profile, user needs to select the 'Change Status To'
  Dropdown in status section as Terminated and then click on Update.
  Agent Termination Page will open. Input the details and click submit
  to terminate the agent.

- **Fields:**

  - Status: dropdown

  - Status Reason: dropdown

  - Status Date: Calendar

  - Effective Date: Calendar

  - Termination Date: Calendar

  - Submit: Button

## 4.3 Agent Goal Setting

- **Purpose:** To set the performance goals for the Agents. The set
  goals should reflect in the agent profile in Agent Portal when the
  agent logins from his ID.

- **Fields:**

  - Agent ID: Text: Option should be given to search & select the agent.

  - Agent Name: Text: Auto-populated

  - Goal Period: From & To Dates from Calendar: To specify timeframe for
    the goal.

  - Target Number of Policies: Textbox: The number of new policies the
    agent aims to sell.

  - Target Premium Collection (₹): Textbox: The total premium amount
    expected to be collected.

  - Product-wise Targets: Textbox: Goals broken down by PLI/RPLI product
    types (e.g., Endowment, Whole Life).

  - Comments: Textbox

  - Submit: button

## 4.4 Agent Commission Management

### 4.4.1 Commission Rate Table View Page

- **Purpose:** To view the Defined commission rates based on multiple
  parameters.

- **Fields:** The page should display the following Commission table:

  - Rate (%): Decimal: Commission percentage

  - Policy Duration (Months): Integer: Duration of policy in months

  - Product Type: Dropdown: PLI, RPLI

  - Product Plan Code: Text: Unique code for the plan

  - Agent Type: Dropdown: e.g., Direct Agent, Field Officer

  - Policy Term (Years): Integer: Total term of the policy

- **Action:**

  - User should only view the table.

### 4.4.2 Commission History Search Page

- **Purpose:** View historical commission data.

- **Search Filters:**

  - Agent ID

  - Policy Number

  - Date Range

  - Product Type

  - Commission Type (First Year, Renewal, Bonus)

- **Result Table:**

  - \| Agent ID \| Policy Number \| Product Type \| Commission Type \|
    Amount \| Status \| Date Processed \|

- **Actions:**

  - Export to Excel/PDF

  - View Detailed Statement

### 4.4.3 Commission Processing Steps:

- **Purpose:** To generate and pay the commission for the Agents.

![](media/image2.png){width="3.638888888888889in"
height="5.458333333333333in"}

#### 4.4.3.1 Run Commission Calculation Batch Jobs

- **Trigger:** Scheduled or Manual

- **Function:** Calculates Commission based on active policies and rate
  table.

#### 4.4.3.2 Trial Statement Generation Batch Job

- **Function:** Generates trial commission statements for review.

- **Output:** Trial Statement per agent and policy.

#### 4.4.3.3 View Trial Statement Page

- **Fields:**

  - Agent ID

  - Policy Number

  - Commission Type

  - Calculated Amount

  - Status (Pending/Approved)

  - Remarks

- **Action:**

  - Filter by Agent, Policy, Circle

  - Export to Excel/PDF

  - Raise Correction

#### 4.4.3.4 Manual Trial Statement Generation Page

- **Fields:**

  - Processing Unit: Dropdown: e.g., IT2.0

  - Statement Format: Dropdown: e.g., Standard

  - Max Statement Due Date: Date: Cut-off for statement

  - Max Transaction Effective Date: Date: Latest transaction date

  - Max Process Date: Date: Latest processing date

  - Statement Date: Date: Date of statement

  - Contract Holder: Text: Name of policyholder

  - Advisor Coordinator: Text: Agent's supervisor

  - Carrier: Text: Insurance carrier

  - Tax Deduction (TDS %): Decimal: Applicable tax deduction

- **Action:**

  - Generate Statement

  - Save Draft

  - Submit for Approval

#### 4.4.3.5 Approving Trial Statement Page

- **Fields:**

  - Agent ID

  - Policy Number

  - Commission Amount

  - Status

  - Remarks

  - Disbursement Option (Full / Partial)

  - Part Disbursement (%) -- Enabled only if selected

- **Actions:**

  - Apply Part Disbursement

  - Approve All Rows in P/U

  - Submit

#### 4.4.3.6 Final Incentive Statement Generation Batch Job

- **Trigger:** Scheduled or Manual

- **Function:** Locks approved trial data and generates final commission
  statements.

#### 4.4.3.7 Final Statements Page

- **Fields:**

  - Agent ID

  - Policy Number

  - Final Commission Amount

  - TDS Deducted

  - Net Payable

  - Payment Status

- **Actions:**

  - View Statement

  - Export PDF/Excel

  - Send to Disbursement

#### 4.4.3.8 Disbursement Details Page

- **Fields:**

  - Agent ID

  - Payment Mode (Cheque / EFT)

  - Cheque Number (if applicable)

  - Bank Name

  - IFSC Code

  - Account Number

  - Payment Date

  - Amount Paid

  - Remarks

- **Actions:**

  - Save

  - Submit

  - Generate Payment File

#### 4.4.3.9 Automatic Disbursement Option

- **Functionality:**

  - If Cheque is selected:

    - Disbursement is marked as complete immediately.

  - If EFT is selected:

    - Payment file is generated and sent to PFMS/Bank.

    - Status updated upon confirmation.

## **5. Test Case**

  --------------------------------------------------------------------------------------------
  **TC     **Functionality**   **Test Case       **Input Data**     **Expected      **Type**
  ID**                         Description**                        Result**        
  -------- ------------------- ----------------- ------------------ --------------- ----------
  TC_001   Agent Onboarding    Create new        Advisor details +  Advisor profile Positive
                               Advisor linked to valid Coordinator  created         
                               existing          ID                 successfully    
                               Coordinator                                          

  TC_002   Agent Onboarding    Create Advisor    Advisor details    Error:          Negative
                               without linking   only               "Coordinator ID 
                               Coordinator                          is mandatory"   

  TC_003   Coordinator         Create new        Coordinator        Coordinator     Positive
           Onboarding          Coordinator with  details + Circle + profile created 
                               valid Circle &    Division           successfully    
                               Division                                             

  TC_004   Coordinator         Create            Coordinator        Error: "Circle  Negative
           Onboarding          Coordinator       details only       assignment      
                               without Circle                       required"       
                               assignment                                           

  TC_005   Dept Employee       Auto-populate     Employee ID from   Profile         Positive
           Onboarding          using valid       HRMS               populated       
                               Employee ID                          correctly       

  TC_006   Dept Employee       Auto-populate     Invalid ID         Error:          Negative
           Onboarding          using invalid                        "Employee ID    
                               Employee ID                          not found"      

  TC_007   Field Officer       Manual entry of   Valid details      Profile created Positive
           Onboarding          details           entered manually   successfully    

  TC_008   Field Officer       Manual entry with Missing Name or ID Error:          Negative
           Onboarding          missing mandatory                    "Mandatory      
                               fields                               fields missing" 

  TC_009   Agent Search        Search by valid   Agent ID = 12345   Agent profile   Positive
                               Agent ID                             displayed       

  TC_010   Agent Search        Search by invalid Agent ID = 99999   Error: "No      Negative
                               Agent ID                             records found"  

  TC_011   PAN Update          Update PAN with   PAN = ABCDE1234F   PAN updated     Positive
                               valid format                         successfully    

  TC_012   PAN Update          Update PAN with   PAN = 12345ABCDE   Error: "Invalid Negative
                               invalid format                       PAN format"     

  TC_013   Status Update       Change status to  Status =           Status updated  Positive
                               Suspended with    Suspended, Reason  successfully    
                               reason            = "Non-compliance"                 

  TC_014   Status Update       Change status     Status =           Error: "Reason  Negative
                               without reason    Suspended, Reason  is mandatory"   
                                                 = blank                            

  TC_015   License Reminder    Generate reminder License expiry =   Reminder        Positive
                               15 days before    15 days ahead      generated       
                               expiry                                               

  TC_016   License Reminder    Generate reminder License expired    No reminder,    Negative
                               after expiry                         license         
                                                                    deactivated     

  TC_017   Commission Rate     Add valid         Rate = 5%, Product Rate saved      Positive
           Setup               commission rate   Type = Endowment   successfully    

  TC_018   Commission Rate     Add rate without  Rate = 5%, Product Error: "Product Negative
           Setup               Product Type      Type = blank       Type required"  

  TC_019   Commission Batch    Execute monthly   Policies exist for Commission      Positive
                               batch with valid  month              calculated      
                               policies                                             

  TC_020   Commission Batch    Execute batch     No policies for    Error: "No data Negative
                               with no policies  month              to process"     

  TC_021   Approve Trial       Approve full      Valid trial        Approval        Positive
           Statement           disbursement      statement          successful      

  TC_022   Approve Trial       Approve without   No selection       Error:          Negative
           Statement           selecting                            "Disbursement   
                               disbursement type                    type required"  

  TC_023   Commission          Auto disbursement Valid EFT details  Commission      Positive
           Disbursement        via EFT                              queued for EFT  

  TC_024   Commission          Disbursement      Missing cheque/EFT Error: "Payment Negative
           Disbursement        without payment   info               details         
                               details                              required"       
  --------------------------------------------------------------------------------------------

## **6. Appendices**

The Following Documents attached below can be used.

![](media/image3.emf)

---

# COMPREHENSIVE COMMISSION AND INCENTIVE MANAGEMENT REQUIREMENTS

## 7. Commission Structure and Rates

### 7.1 PLI - First Year Incentive

| Policy Type | Incentive Rate (First Year Premium) |
|-------------|-------------------------------------|
| Other than AEA - Term <= 15 years | 4% |
| Other than AEA - 15 < Term <= 25 years | 10% |
| Other than AEA - Term > 25 years | 20% |
| AEA - Term <= 15 years | 5% |
| AEA - Term > 15 years | 7% |

### 7.2 RPLI - First Year Incentive

- 10% of the first-year premium income for RPLI policies procured by the sales force.

### 7.3 Renewal Incentive

| Policy Procurement Period | Eligible Agents | Renewal Incentive Rate |
|---------------------------|-----------------|------------------------|
| 01.10.2009 - 31.03.2017 | Field Officers, Direct Agents | 1% of renewal premium |
| On or after 01.07.2020 | All Agent Types | 1% of renewal premium |
| PLI (Cash Policies) | All Agent Types | 1% of renewal premium |
| RPLI Policies | All Agent Types | 2.5% of renewal premium |

**Note:** No renewal incentive for pay policy.

### 7.4 Monitoring Staff - Procurement Incentive

| Monitoring Role | Rate on New Business Premium |
|-----------------|------------------------------|
| Development Officer (for DA/FO under him) | 0.8% |
| Sub-Divisional Head (for GDS under him) | 0.6% |
| Mail Overseer (for GDS under him) | 0.2% |
| Sub-Divisional Head (for DE under him) | 0.8% |
| ASP(HQ)/OS or delegated Dy. SP (For DE under him) | 0.8% |
| Divisional Head (for all sales force) | 0.2% |

**Note:** No renewal incentive is payable to monitoring staff.

## 8. Eligibility and Calculation Rules

1. **Commission Eligibility:** Commission will be calculated only after the policy is accepted and free-look period (15 days) is completed.

2. **Premium Realization:** Incentive shall be payable only on realized premiums. In case of online payment T+5.

3. **Monthly Processing:** Incentives shall be computed once every month.

4. **Free-Look Period Exclusion:** All computations will consider policies only after acceptance and completion of the 15-day free-look period.

## 9. Payment Flow

The commission and incentive payment shall follow this workflow:

1. System computes commission for eligible policies post free-look.
2. Approval statements are generated for verification by Divisional Head.
3. Upon approval, payment files are generated for disbursement through Finacle.
4. All transactions shall be logged for audit with timestamps and user IDs.

## 10. Taxation

1. **TDS Deduction:** TDS @ 2% shall be deducted at source from the payable incentive.
2. **GST Payment:** GST @ 18% shall be paid by the Department on a reverse charge basis (RCM).
3. **Reporting:** The system shall generate monthly TDS and GST liability reports.

## 11. Functional Requirements (Extended)

| ID | Functionality | Requirement |
|----|---------------|-------------|
| FR-01 | Configure Commission and Incentive Rates | System shall allow configuration of commission rates by policy type, term, and agent category |
| FR-02 | Capture Acceptance and Free-Look Data | System shall capture policy acceptance date and free-look completion |
| FR-03 | Compute First Year and Renewal Incentives | System shall calculate incentives based on defined rate tables |
| FR-04 | Apply Procurement Incentive Rules for Monitoring Staff | System shall allocate procurement incentives to monitoring hierarchy |
| FR-05 | Generate Approval Statements | System shall generate trial statements for approval |
| FR-06 | Integrate with Finacle for payment disbursement | System shall generate payment files for Finacle integration |
| FR-07 | Maintain audit trail of all computations and payments | System shall log all transactions with timestamps and user IDs |
| FR-08 | Generate reports and dashboards for management review | System shall provide MIS reports for management |

## 12. Reports

1. **Agent-wise Incentive Summary** (First Year and renewal)
2. **Circle wise/Division-wise Incentive Summary**
3. **Monitoring Staff Procurement Incentive Register**
4. **TDS and GST Summary Report**
5. **Pending Approval / Disbursement Report**
6. **Policy Wise Incentive Report**
7. **Agent category wise incentive Report** (First Year and Renewal)

## 13. Non-Functional Requirements

1. **Performance:** System shall process monthly commission batches within prescribed time limits.
2. **Audit:** All transactions and computations shall be auditable and traceable.
3. **Access Control:** Role-based access control shall be implemented.
4. **Security:** Data shall be encrypted in transit and at rest.
5. **Retention:** System shall maintain historical records for minimum 10 years.

## 14. Sample Test Cases (Extended)

| ID | Description |
|----|-------------|
| TC-01 | Verify correct incentive rate for AEA vs non-AEA policies |
| TC-02 | Validate free-look exclusion logic |
| TC-03 | Check correct allocation to monitoring staff roles |
| TC-04 | Verify Finacle payment integration for approved statements |
| TC-05 | Confirm TDS deduction and GST computation correctness |

## 15. Integration Specifications

### 15.1 Finacle Integration

**Purpose:** Disbursement of approved commission and incentive payments

**Integration Type:** File-based/API (to be confirmed with CEPT team)

**Integration Flow:**
1. System generates payment file post-approval
2. File format: CSV/XML with following fields:
   - Beneficiary Account Number
   - Beneficiary Name
   - Amount (Net of TDS)
   - Transaction Reference Number
   - Payment Date
   - Remark
3. File encryption before transfer
4. Transfer Protocol: SFTP to designated Finacle folder
5. Frequency: Twice in a week (for approved batches)
6. Acknowledgment: Finacle to provide status file (Success/Failure)
7. Reconciliation: Daily reconciliation of payment status

**Error Handling:**
- Failed payments to be flagged and notified to CEPT
- Retry mechanism for network failures (max 3 attempts)
- Manual intervention queue for persistent failures

### 15.2 Policy Management System Integration

**Purpose:** Fetch policy acceptance, free-look completion, and premium realization data

**Data Required:**
- Policy Number
- Policy Type (PLI/RPLI, AEA/Non-AEA)
- Policy Term
- Acceptance Date
- Free-look Period Completion Date
- Premium Amount (First Year/Renewal)
- Premium Due Date
- Premium Realization Date
- Premium Realization Status
- Payment Mode (Online/Offline)
- Agent Code
- Monitoring Staff Hierarchy

**Integration Frequency:** Real-time API calls or Daily batch file

**Data Validation:** Mandatory field checks, date validations, amount validations

### 15.3 Agent Master System Integration

**Purpose:** Validate agent eligibility and status

**Data Required:**
- Agent Code
- Agent Name
- Agent Type (DE/Retired FO/Direct Agent/GDS)
- Status (Active/Inactive/Suspended/Terminated)
- Bank Account Details
- PAN Number
- Joining Date
- Reporting Officer Code (for monitoring staff linkage)
- GST Number

**Integration Frequency:** Daily sync for master data updates

### 15.4 Employee Master Integration

**Purpose:** Validate monitoring staff and departmental employees

**Data Required:**
- Employee Code
- Employee Name
- Designation
- Posting Location (Division/Sub-Division)
- Active Status
- Transfer/Promotion dates
- Bank Account Details

**Integration Frequency:** Daily sync

## 16. Detailed Business Rules and Calculations

### 16.1 Agent Eligibility Criteria

**Active Status:**
- Agent must be in 'Active' status on the date of policy acceptance
- If agent status changes to Inactive/Suspended/Terminated after policy acceptance but before commission processing, commission shall still be payable

**Minimum Service Period:**
- No minimum service period required for first-year commission
- For renewal commission, agent must have been active at the time of policy procurement

**Multiple Agents:**
- Only the primary procuring agent shall receive commission
- Joint procurement scenarios not applicable

### 16.2 Premium Realization Rules

**Offline Payment (Cash/Cheque/DD):**
- Commission eligible once premium is credited to Department account
- Realization date = Credit date

**Online Payment (Net Banking/UPI/Card):**
- Commission eligible after T+5 days from transaction date
- T+5 buffer to account for payment gateway settlement and chargeback window
- Realization date = Transaction date + 5 working days

**Premium Bounce/Reversal:**
- If premium payment fails/bounces after commission processing:
  - Commission to be reversed in next cycle
  - Negative entry in agent's commission statement
  - Net payable adjusted in subsequent month

### 16.3 Pro-rata Calculations

**Mid-month Policy Acceptance:**
- No pro-rata adjustment
- Full commission payable in the month of free-look completion

**Mid-month Agent Status Change:**
- If agent becomes inactive mid-month, commission for policies accepted before status change shall be paid
- If agent joins mid-month, commission for policies accepted after joining date shall be paid

**Mid-month Transfer of Monitoring Staff:**
- Commission allocated based on posting on the date of policy acceptance
- If monitoring staff transfers mid-month, new incumbent gets commission for policies accepted after transfer date

### 16.4 Monitoring Staff Hierarchy Validation

**Hierarchy Rules:**
- System shall validate reporting hierarchy at the time of commission calculation
- Commission allocated to monitoring staff based on hierarchy as on policy acceptance date

**Multiple Levels:**
- Development Officer receives 0.8% for policies by DA/FO reporting to them
- Sub-Divisional Head receives 0.8% (depending on GDS/DE reporting)
- Mail Overseer receives 0.8% for GDS under supervision
- ASP(HQ)/OS receives 0.8% for DE under supervision
- Divisional Head receives 0.2% for all sales force in division

**Hierarchy Change:**
- If agent reporting changes mid-month, monitoring staff as on policy acceptance date receives commission
- No retrospective adjustment for hierarchy changes

### 16.5 Commission Clawback Scenarios

**Triggers for Clawback:**
1. Premium bounce/reversal
2. Data entry error correction
3. Fraudulent policy detection
4. Agent disciplinary action resulting in policy invalidation

**Clawback Process:**
1. Identify affected policies and commission amounts
2. Generate clawback statement for affected agents
3. Adjust in next commission cycle as negative entry
4. If insufficient future commission, initiate recovery through establishment
5. Maintain clawback register for audit

## 17. Exception Handling and Error Scenarios

### 17.1 Data Validation Failures

**Missing Agent Details:**
- Flag policy for manual review
- Generate exception report for admin
- Hold commission until resolution
- SLA: To be decided by CEPT

**Invalid Policy Data:**
- Premium amount = 0 or negative
- Policy term outside defined range
- Missing acceptance date
- Action: Reject calculation, log error, notify source system

**Hierarchy Mismatch:**
- Agent code not mapped to monitoring staff
- Monitoring staff not found in employee master
- Action: Allocate agent commission, hold monitoring staff commission, generate exception report

### 17.2 Premium Realization Issues

**Premium Not Realized:**
- Commission calculation on hold
- Monitor for 90 days
- If not realized within 90 days, mark policy as "Premium Pending"
- No commission payable until realization

**Partial Realization:**
- Calculate commission on realized portion
- Balance commission on hold until full realization
- Track partial payments and cumulative commission

### 17.3 Integration Failures

**Finacle Integration Failure:**
- Retry 3 times with 1-hour interval
- If all retries fail, move to manual queue
- Send email alert to CEPT and Finance team
- Generate failure report with error details

**Source System Data Feed Failure:**
- Use previous day's data for critical reports
- Flag calculations as "Provisional"
- Reconcile and adjust once feed is restored

### 17.4 Calculation Errors

**Rate Configuration Error:**
- System should validate rate configuration before processing
- If error detected mid-processing, halt batch
- Notify admin team
- Re-run after correction

**Arithmetic Overflow/Underflow:**
- Implement boundary checks (min: Re. 1, max: Rs. 50,000 per policy)
- Log exceptions for manual review

## 18. User Interface Requirements

### 18.1 Dashboard - Management View

**KPIs to Display:**
- Total Commission Payable (Current Month)
- Total Policies Processed
- Total Agents Eligible
- Average Commission per Agent
- Circle-wise/Division-wise Breakdown
- Trend Analysis (Last 6 months)
- Top 10 Performing Agents
- Exception Count (Pending Approvals, Errors)

**Filters:**
- Date Range
- Circle/Division/Sub-Division
- Agent Type
- Policy Type
- Commission Type (First Year/Renewal)

**Visualization:**
- Bar charts for circle-wise comparison
- Line charts for trend analysis
- Pie charts for policy type distribution
- Data tables with drill-down capability

### 18.2 Commission Calculation Screen

**User Role:** Commission Processing Officer

**Features:**
1. Select calculation period (month-year)
2. Select scope (All/Specific Circle/Division)
3. Run calculation button
4. Progress indicator during processing
5. Summary of calculated records
6. View detailed calculations (agent-wise)
7. Export to Excel functionality
8. Error/Exception log viewer

**Validations:**
- Cannot process same period twice (unless previous calculation is cancelled)
- Cannot process future periods
- Must complete previous period before processing next

### 18.3 Approval Workflow Screen

**User Role:** Approving Authority (Divisional Head)

**Features:**
1. View pending approval queue
2. Filter by division/agent type/amount range
3. View detailed calculation breakdown for selected entry
4. Bulk approve functionality
5. Individual approve/reject with remarks
6. Compare with previous month data
7. Audit trail view

**Approval Levels:**
- Level 1: Divisional Head

### 18.4 Agent Self-Service Portal

**Features for Agents:**
1. View own commission statements (month-wise)
2. View policy-wise commission breakdown
3. Download commission certificate
4. View payment status
5. Raise query/dispute
6. View TDS certificate

**Security:**
- Login using Agent Code and OTP
- View only own data
- No edit permissions

### 18.5 Search and Filter Requirements

**Global Search:**
- Search by Agent Code, Agent Name, Policy Number, Transaction ID, PAN, Mobile
- Auto-suggest functionality
- Recent searches history

**Advanced Filters:**
- Date range (from-to)
- Amount range
- Status (Pending/Approved/Paid/Rejected/On Hold)
- Circle/Division/Sub-Division
- Agent Type
- Policy Type
- Save filter presets for frequently used combinations

### 18.6 Export Capabilities

**Supported Formats:**
- Excel (.xlsx)
- PDF (for statements and certificates)
- CSV (for bulk data)

**Export Options:**
- Current page data
- Selected records
- All filtered records (with warning for large datasets)
- Include/exclude filters and summary

**Security:**
- Watermark on exported PDFs
- Log all export activities with user ID and timestamp
- Limit bulk exports (max 50,000 records per export)

## 19. Security and Access Control

### 19.1 Role Definitions and Permissions Matrix

| Role | Permissions |
|------|-------------|
| Divisional Head | Run calculations, View all agent data, Generate reports, View division data, Approve Incentive/Commission, NO calculation rights |
| Circle Head | View circle data, Access all reports |
| Regional Head | View region data, Access all reports |
| PLI Directorate | View circle data, Access all reports |
| CEPT | User management, System configuration, Troubleshooting, System configuration rights (rate changes) |
| Audit | Read-only access to all data, Audit trail access, Report generation, NO edit/approve rights |
| Agent (Self-Service) | View own commission data, Download own statements, NO access to other agents' data |

### 19.2 Maker-Checker Requirements

**Critical Operations Requiring Maker-Checker:**
1. Rate configuration changes
2. Manual commission adjustments
3. Clawback entries
4. Payment file generation
5. User role assignments

**Workflow:**
- Maker initiates the transaction
- System logs maker ID and timestamp
- Transaction moves to checker queue
- Checker reviews and approves/rejects
- System logs checker ID and timestamp
- On rejection, transaction returns to maker with comments

**Segregation Rules:**
- Same user cannot be maker and checker for same transaction
- Checker must be of higher or equal grade than maker

### 19.3 Data Privacy and Confidentiality

**PII (Personally Identifiable Information) Protection:**
- Agent bank account numbers masked (show last 4 digits only)
- PAN numbers masked (show last 4 characters only)
- Phone numbers and email addresses encrypted in database
- Access to full PII restricted to authorized roles only

**Data Masking Rules:**
- Non-privileged users see masked data in UI
- Export files for non-privileged users contain masked data
- Audit logs capture access to sensitive data

**Compliance:**
- Adherence to IT Act and data protection regulations
- Data retention as per government guidelines
- Secure disposal of data after retention period

### 19.4 Session Management

**Session Security:**
- Session timeout: 5 minutes of inactivity
- Concurrent session prevention (one active session per user)
- Force logout on role change
- IP address logging for each session

**Password Policy:**
- Minimum 8 characters (1 uppercase, 1 lowercase, 1 number, 1 special character)
- Password expiry: 180 days
- Password history: Cannot reuse last 5 passwords
- Account lockout: 3 failed login attempts (unlock after 30 minutes or by admin)

### 19.5 Audit and Compliance

**Audit Trail Capture:**
All user actions logged with:
- User ID
- Timestamp
- Action performed
- Old value and new value (for updates)
- IP address
- Session ID

**Audit Trail Retention:**
- Minimum 10 years as per statutory requirements
- Archive older data to separate storage
- Search and retrieval capability for audit data

**Compliance Reporting:**
- Generate audit reports for internal audit
- CAG audit support (data extraction and reconciliation)
- RTI compliance (mechanism to extract specific data)

## 20. Performance Requirements

### 20.1 Response Time Requirements

| Operation | Expected Response Time |
|-----------|------------------------|
| User Login | < 2 seconds |
| Dashboard Load | < 3 seconds |
| Search Results | < 2 seconds |
| Report Generation (< 1000 records) | < 5 seconds |
| Report Generation (< 10,000 records) | < 30 seconds |
| Commission Calculation (monthly batch) | < 2 hours |
| Payment File Generation | < 15 minutes |
| Export to Excel (< 5000 records) | < 10 seconds |

### 20.2 Batch Processing Windows

**Daily Tasks:** To be decided by CEPT

### 20.3 Concurrent User Requirements

**Peak Load Handling:** To be decided by CEPT

**Stress Testing:**
- System should be tested for 150% of expected peak load

### 20.4 Availability Requirements

**Uptime SLA:**
- System availability: 99.5% (excluding planned maintenance)
- Advance notification (48 hours) for planned outages

**Disaster Recovery:**
To be decided by CEPT

## 21. Technical Specifications

**Application Layer:**
- Web-based application (Browser-based access)
- Responsive design (desktop and tablet support)
- Supported browsers: Chrome, Firefox, Edge (latest 2 versions)

## 22. Abbreviations and Business Terms

### 22.1 Abbreviations

| Term | Full Form |
|------|-----------|
| AEA | Anticipated Endowment Assurance |
| PLI | Postal Life Insurance |
| RPLI | Rural Postal Life Insurance |
| FO | Field Officer |
| DA | Direct Agent |
| DE | Departmental Employee |
| GDS | Gramin Dak Sevak |
| ASP | Assistant Superintendent of Post Offices |
| HQ | Headquarters |
| OS | Office Superintendent |
| Dy. SP | Deputy Superintendent |
| DO | Development Officer |
| SDH | Sub-Divisional Head |
| DH | Divisional Head |
| CBS | Core Banking System |
| TDS | Tax Deducted at Source |
| GST | Goods and Services Tax |
| RCM | Reverse Charge Mechanism |
| PAN | Permanent Account Number |
| UAT | User Acceptance Testing |
| SLA | Service Level Agreement |
| RTO | Recovery Time Objective |
| RPO | Recovery Point Objective |
| PII | Personally Identifiable Information |

### 22.2 Business Terms

**Free-Look Period:**
- 15-day window from policy issuance during which policyholder can cancel without penalty
- Commission not payable until free-look period completion

**Premium Realization:**
- Actual receipt and confirmation of premium payment in Department's account
- Commission calculated only on realized premiums

**First Year Premium:**
- Total premium collected in the first policy year
- Basis for first-year incentive calculation

**Renewal Premium:**
- Premium collected in subsequent policy years (2nd year onwards)
- Basis for renewal incentive calculation

**Cash Policy:**
- PLI policies where premium is paid other than pay
- Separate renewal incentive rate applicable

**Policy Term:**
- Duration of the policy in years
- Determines applicable incentive rate

**Procurement:**
- Enrolling a new policy
- Used in context of monitoring staff incentive

**Clawback:**
- Recovery of commission already paid due to policy cancellation/reversal

**T+5:**
- Transaction date plus 5 working days
- Settlement period for online payments before commission eligibility

### 22.3 System Terms

**Batch Processing:**
Automated execution of commission calculation for all eligible policies in a defined period

**Approval Workflow:**
Sequential process of review and authorization before payment

**Maker-Checker:**
Dual control mechanism where one user initiates and another approves

**Audit Trail:**
Chronological record of all system activities and changes

**Role-Based Access Control (RBAC):**
Permission management based on user roles

**Data Masking:**
Hiding sensitive data from unauthorized users

## 23. Assumptions and Dependencies

### 23.1 Assumptions

**1. Data Availability:**
- Policy acceptance data is available in source system within 1 day of acceptance
- Premium realization data is updated real-time or within T+1 day
- Agent master data is maintained accurately and updated promptly

**2. Business Process:**
- Approval authorities are available within defined SLA for approvals
- Bank account details of agents are validated and active
- PAN details of agents are accurate for TDS deduction

**3. User Readiness:**
- Users are trained on the new system before go-live
- User manuals and SOPs are prepared and distributed
- Help desk support is available post go-live

**4. Regulatory:**
- Current TDS and GST rates remain unchanged during implementation
- No major policy changes from Department during implementation

### 23.2 Dependencies

**External Dependencies:**

**1. Finacle System:**
- Availability of Finacle API/file interface for payment processing
- Timely processing of payment files by Finacle team
- Reconciliation support from Finacle team

**2. Policy Management System:**
- Availability of real-time/daily data feed
- Data quality and accuracy from source
- Technical support for integration issues

**3. Network and Infrastructure:** To be decided by CEPT

**4. Third-Party Services:**
- SMS gateway for OTP and alerts
- Email server for notifications

**Internal Dependencies:**

**1. Data Readiness:**
- Historical agent and policy data cleanup and migration
- Master data validation and correction
- Mapping of agent hierarchy and monitoring staff

**2. Business Readiness:**
- Approval of commission rates by competent authority
- Finalization of approval workflow and authorities
- Resolution of business rule ambiguities

**3. IT Infrastructure:** To be decided by CEPT

### 23.3 Risks and Mitigation

| Risk | Impact | Mitigation |
|------|--------|------------|
| Data quality issues in source systems | High | Conduct data profiling, implement validation rules, manual verification for exceptions |
| Integration delays with Finacle | High | Early engagement with Finacle team, parallel testing, fallback to manual process if needed |
| User resistance to new system | Medium | Comprehensive training, change champions, continuous support |
| Performance issues with large data volumes | Medium | Performance testing, optimization, scalable architecture |
| Regulatory changes during implementation | Low | Flexible system design, configurable rules, quick response mechanism |

## 24. Success Criteria for UAT Acceptance

### 24.1 Functional Acceptance Criteria

**1. Calculation Accuracy:**
- 100% accuracy in commission calculation for all policy types
- Correct application of rates based on policy term and type
- Accurate TDS and GST computation

**2. Workflow Completion:**
- End-to-end workflow from calculation to payment successful
- All approval levels functioning correctly
- Payment file generation and Finacle integration successful

**3. Exception Handling:**
- All defined exception scenarios handled correctly
- Error messages are clear and actionable
- Exception reports generated accurately

**4. Reports:**
- All 7 defined reports generate correctly with accurate data
- Export functionality works for all reports
- Report filters and drill-downs functioning

**5. User Interface:**
- All screens load without errors
- Navigation is intuitive and as per design
- Search and filter functionality works correctly

### 24.2 Non-Functional Acceptance Criteria

**1. Performance:**
- Response times meet defined SLAs for all operations
- Monthly batch processing completes within 2 hours

**2. Security:**
- All security controls implemented and tested
- Penetration testing completed with no critical vulnerabilities
- Data masking working correctly for non-privileged users

**3. Audit:**
- Audit trail captures all defined events
- Audit reports generate correctly
- Historical data retrievable for 10 years

**4. Integration:**
- All integrations tested and functioning
- Error handling for integration failures working
- Reconciliation processes successful

### 24.3 UAT Sign-off Checklist

- [ ] All test cases executed successfully (minimum 95% pass rate)
- [ ] Critical and high-priority defects resolved
- [ ] Performance benchmarks achieved
- [ ] Security assessment completed
- [ ] User training completed
- [ ] User manual and documentation delivered
- [ ] Data migration completed and validated
- [ ] Reconciliation with legacy system successful
- [ ] Go-live checklist prepared
- [ ] Rollback plan documented
- [ ] Support team trained and ready
- [ ] Business sign-off obtained from all stakeholders

## 25. Reconciliation and Validation

### 25.1 Daily Reconciliation

**Policy Data Reconciliation:**
- Count of policies received from source system
- Count of policies processed for commission
- Identify unprocessed policies with reasons

**Payment Reconciliation:**
- Payment files sent to Finacle
- Payment confirmations received
- Failed payments
- Pending payments

### 25.2 Monthly Reconciliation

**Commission Calculation Reconciliation:**
- Total commission calculated vs. approved
- Total commission approved vs. paid
- Pending approvals
- Rejected commissions with reasons

**Financial Reconciliation:**
- Total premium collected (from policy system)
- Total commission calculated (from commission system)
- Commission as % of premium (variance analysis)
- TDS and GST reconciliation with finance books

### 25.3 Annual Reconciliation

**Agent-wise Annual Statement:**
- Total policies procured
- Total commission earned
- Total TDS deducted
- Net commission paid
- Comparison with previous year

**Monitoring Staff Annual Statement:**
- Total procurement incentive earned
- Division/Circle wise breakup
- Comparison with targets

**Audit Reconciliation:**
- Total commission paid vs. budget
- Policy-wise commission register
- Exception register
- Clawback register

## 26. Compliance and Regulatory Requirements

### 26.1 Statutory Compliance

**Income Tax Compliance:**
- TDS deduction as per Income Tax Act
- Form 16A generation for agents
- Quarterly TDS return filing (Form 24Q/26Q)
- Annual TDS reconciliation

**GST Compliance:**
- GST payment under RCM
- Monthly GST return filing
- GST reconciliation with books
- GST audit trail maintenance

**Government Accounting:**
- Integration with government accounting system
- Month-end and year-end closures
- CAG audit requirements

## 27. Future Enhancements (Out of Scope for Current Release)

### 27.1 Potential Future Features

**Agent Performance Analytics:**
- Predictive analytics for agent performance
- Identification of high-performing agents
- Territory performance heat maps

**AI/ML Integration:**
- Anomaly detection in commission calculations
- Fraud detection in policy procurement
- Automated exception resolution
- Chatbot for agent queries

**Advanced Reporting:**
- Self-service BI tool integration
- Customizable dashboards
- Predictive revenue forecasting
- What-if scenario analysis

**Cloud Migration:** (if applicable)

## 28. Appendices (Extended)

### 28.1 Sample Commission Calculation Examples

**Example 1: PLI Non-AEA Policy (Term 20 years)**
- Policy Type: PLI Whole Life
- Policy Term: 20 years
- First Year Premium: Rs. 12,000
- Agent Type: Direct Agent
- Calculation:
  - Incentive Rate: 10% (15 < Term <= 25 years)
  - Gross Commission: 12,000 x 10% = Rs. 1,200
  - TDS @ 2%: Rs. 24
  - Net Payable: Rs. 1,176

**Example 2: RPLI Policy**
- Policy Type: RPLI
- First Year Premium: Rs. 5,000
- Agent Type: GDS
- Calculation:
  - Incentive Rate: 10%
  - Gross Commission: 5,000 x 10% = Rs. 500
  - TDS @ 2%: Rs. 10
  - Net Payable: Rs. 490
- Monitoring Staff:
  - Sub-Divisional Head: 5,000 x 0.6% = Rs. 30
  - Mail Overseer: 5,000 x 0.2% = Rs. 10
  - Divisional Head: 5,000 x 0.2% = Rs. 10

**Example 3: Renewal Incentive (RPLI)**
- Policy Procured: 15.08.2021
- Renewal Premium: Rs. 4,000
- Agent Type: Field Officer
- Calculation:
  - Renewal Incentive Rate: 2.5%
  - Gross Commission: 4,000 x 2.5% = Rs. 100
  - TDS @ 2%: Rs. 2
  - Net Payable: Rs. 98

**Example 4: Multiple Policies for Same Agent**
- Agent Code: DA12345
- Policy 1: PLI, Term 12 years, Premium Rs. 10,000 -> 4% = Rs. 400
- Policy 2: RPLI, Premium Rs. 8,000 -> 10% = Rs. 800
- Policy 3: PLI, Term 30 years, Premium Rs. 15,000 -> 20% = Rs. 3,000
- Total Gross Commission: Rs. 4,200
- TDS @ 2%: Rs. 84
- Net Payable: Rs. 4,116

### 28.2 Sample Reports Format

**Agent-wise Incentive Summary Report (Including Monitoring staff)**

Report Period: April 2025 (01-04-2025 to 30.04.2025)
Generated On: 05-May-2025 10:30 AM

| Agent Code | Agent Name | Agent Type | No. of Policies | Total Premium | Gross Commission | TDS | Net Payable |
|------------|------------|------------|-----------------|---------------|------------------|-----|-------------|
| DA001 | Rajesh Kumar | Direct Agent | 5 | 45,000 | 5,200 | 104 | 5,096 |
| FO023 | Sunita Sharma | Field Officer | 3 | 2,80,000 | 2,800 | 56 | 2,744 |
| | | Total | 8 | 3,25,000 | | | |

**TDS Summary Report (Including Monitoring staff)**

Month: April 2025

| Category | Gross Commission | TDS @ 2% | Net Payable |
|----------|------------------|----------|-------------|
| Direct Agents | 85,000 | 1,700 | 83,300 |
| Departmental Employees | 35,000 | 700 | 34,300 |
| Field Officers | 5,000 | 100 | 4,900 |
| Total | 1,25,000 | 2,500 | 1,22,500 |

**Income Tax and GST Report:**

| Name Of Division | Agent Name | Agent Code | PAN No | Incentive Amount (Gross) | Income Tax | Disbursed Incentive Amount | SGST/UGST | CGST |
|-----------------|------------|------------|--------|--------------------------|------------|---------------------------|-----------|------|
| | | | | | | | | |
