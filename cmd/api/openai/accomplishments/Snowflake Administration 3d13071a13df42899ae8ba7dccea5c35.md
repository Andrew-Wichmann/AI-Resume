# Snowflake Administration

Tags: Cost Optimization, Support
Date: October 1, 2021 → February 1, 2024

# Executive Summary

- Work with other admins, to provision user accounts and resources for client teams with a 24 hour turn over
- Terraform module feature development to allow teams to monitor their warehouses and (most of) the warehouses that their role users use
- Designed and performed quarterly audits of admin accounts to ensure compliance with security standards

# Details

In late 2021, as Mintel was integrating more and more with Snowflake, our primary Snowflake support team in the UK wanted to anoint a US-based admin to offer Snowflake assistance for the time zone. I was chosen and granted account and security admin roles on all three of our Snowflake accounts and owner role for the Gitlab group. At first, this came with not a lot more responsibility. The one service I did provide was verifying SQL queries of role users for developers who didn’t have the MONITOR privilege on the right warehouses. Eventually, I took it upon myself to edit the terraform modules that provision role users to accept a list of roles allowed to monitor the role user warehouse. This MR was accepted and my teammates were then able to verify the SQL queries of their ETL pipelines by themselves.

In 2022, the primary Snowflake support team dissolved and not a single author of the original team remained at the company. The new team leader started holding monthly catchups to stabilize and synchronize. My role during this was to produce documentation on the responsibilities of a snowflake admin. I formalized the rules of 24 hour turn over for response to support tickets.

In 2023, when setting up my teammate as the 2nd US-based Snowflake admin, I noticed that the roles granted to the UK and China-based admins were inconsistent. When I raised this to the new team leader, he asked if I could document it. I created a Google Document and accompanying Sheet that outlined my observations, recommendations, and security infringements. Some interesting findings:

- Some former employee accounts weren’t deactivated (last I checked, mine wasn’t)
- Two admins weren’t using 2FA
- The root user was being used for maintenance instead of individual admin users