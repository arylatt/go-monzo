package monzo

// Images (eg. receipts) can be attached to transactions by uploading these via the attachment API.
// Once an attachment is registered against a transaction, the image will be shown in the transaction detail screen within the Monzo app.
//
// There are two options for attaching images to transactions - either Monzo can host the image, or remote images can be displayed.
type AttachmentsService service
