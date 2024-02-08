# Cost Saving Initiatives

Tags: Cost Optimization, Data Ingestion
Date: May 22, 2023 → June 12, 2023

# Executive Summary

- “Sunsetting” a defunct data vendor stack has saved at least ~$450/mo.
- Lifecycling S3 objects in our three biggest buckets to Standard-IA and Glacier-IR has saved ~$3.7k/mo.
- There is still a potential $380-$400/mo to save from migrating clients off of a legacy site.
- New AWS resources left running that relate to this work will cost ~$150/mo to run

# Details

With a lull in management’s requests for new features, our team was able to address some low-hanging cost saving projects.

One of the projects was removing all the AWS resources related to a data vendor that was replaced. This project was difficult because the terraform that provisioned everything was highly coupled to all the other vendors. Trying to change a single vendor would require synchronizing the configuration between terraform code and AWS for all the other vendors. To avoid this, we decided to cut a corner and target apply just the resources relevant to the defunct vendor. Using resource tags, we were able to measure the $450/mo reduction in our AWS bill.

Another project came to us from our architect who noticed that we weren’t making use of any S3 lifecycle policies. Using S3 Storage Lens we were able to confirm that objects in our buckets were primarily being accessed within the first month of being created. Therefore, we implemented the most aggressive policy of lifecycling objects to Standard-IA after 30 days and then to Glacier-IR after 90 days. The initial migration of objects to these storage classes was very expensive due to the large number of objects needing to be lifecycled. IIRC, the cost was somewhere around $11k. Thankfully, however, we were able to get ahead of any panicked emails from accounting by showing that the newly lifecycled objects would save ~$3.7k/mo.

During our investigation of old S3 buckets, we set up a logging pipeline to ship S3 access logs to Loki and gain new insight into exactly what IAM roles are requesting what objects from what buckets and how often. One interesting pattern we noticed was access to a bucket for the aforementioned data vendor. After making some calls and emails, we discovered that an older site’s image server was still using this bucket to serve images, but that any clients that have access to that site should already have access to the new site. Therefore, we asked the maintainers of the old site to remove the search filter option to filter for those images, and we asked account managers to start moving clients off of this old site. After that, we would be clear to delete this bucket.

It was fitting that during a cost saving initiative that we should be accountable for the new resources intentionally left running. Again, using our resource tagging, we were able to communicate that the logging pipeline and S3 storage lens left running would cost about $150. We weren’t able to confirm the exact total because we had no way of estimating the cost to Loki for adding a new log stream.