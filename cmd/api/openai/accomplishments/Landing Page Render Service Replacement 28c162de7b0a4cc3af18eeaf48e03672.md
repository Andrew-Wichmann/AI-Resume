# Landing Page Render Service Replacement

Tags: Data Ingestion, Design
Date: November 19, 2023 → December 20, 2023

# Executive Summary

- The sales team was getting embarrassed by the amount of images of landing pages with errors showing up in search results during sales pitches, so we decided to upgrade the service that renders these images.
- 4 different render services were compared by rendering the same landing page with each service and comparing the results side-by-side. 2 SaaS solutions and 2 open source solutions were compared. [The open source “ws-screenshot” service](https://github.com/elestio/ws-screenshot) was selected due to it’s superior rendering results.
- The new service was deployed
- An audit of 500 images before the deployment and 500 after the deployment was conducted and, using a custom grading rubric, we determined that image quality improved ~18%