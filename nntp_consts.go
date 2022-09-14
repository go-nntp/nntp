package nntp

const (
	// Connection

	/**
	 * 'Server ready - posting allowed' (RFC977)
	 */
	ResponseCodeReadyPostingAllowed = 200

	/**
	 * 'Server ready - no posting allowed' (RFC977)
	 */
	ResponseCodeReadyPostingProhibited = 201

	/**
	 * 'Closing connection - goodbye!' (RFC977)
	 */
	ResponseCodeDisconnectingRequested = 205 ///// goodbye

	/**
	 * 'Service discontinued' (RFC977)
	 */
	ResponseCodeDisconnectingForced = 400 ///// unavailable / discontinued

	/**
	 * 'Slave status noted' (RFC977)
	 */
	ResponseCodeSlaveRecognized = 202

	// Common errors

	/**
	 * 'Command not recognized' (RFC977)
	 */
	ResponseCodeUnknownCommand = 500

	/**
	 * 'Command syntax error' (RFC977)
	 */
	ResponseCodeSyntaxError = 501

	/**
	 * 'Access restriction or permission denied' (RFC977)
	 */
	ResponseCodeNotPermitted = 502

	/**
	 * 'Program fault - command not performed' (RFC977)
	 */
	ResponseCodeNotSupported = 503

	// Group selection

	/**
	 * 'Group selected' (RFC977)
	 */
	ResponseCodeGroupSelected = 211

	/**
	 * 'No such news group' (RFC977)
	 */
	ResponseCodeNoSuchGroup = 411

	// Article retrieval

	/**
	 * 'Article retrieved - head and body follow' (RFC977)
	 */
	ResponseCodeArticleFollows = 220

	/**
	 * 'Article retrieved - head follows' (RFC977)
	 */
	ResponseCodeHeadFollows = 221

	/**
	 * 'Article retrieved - body follows' (RFC977)
	 */
	ResponseCodeBodyFollows = 222

	/**
	 * 'Article retrieved - request text separately' (RFC977)
	 */
	ResponseCodeArticleSelected = 223

	/**
	 * 'No newsgroup has been selected' (RFC977)
	 */
	ResponseCodeNoGroupSelected = 412

	/**
	 * 'No current article has been selected' (RFC977)
	 */
	ResponseCodeNoArticleSelected = 420

	/**
	 * 'No next article in this group' (RFC977)
	 */
	ResponseCodeNoNextArticle = 421

	/**
	 * 'No previous article in this group' (RFC977)
	 */
	ResponseCodeNoPreviousArticle = 422

	/**
	 * 'No such article number in this group' (RFC977)
	 */
	ResponseCodeNoSuchArticleNumber = 423

	/**
	 * 'No such article found' (RFC977)
	 */
	ResponseCodeNoSuchArticleId = 430

	// Transferring

	/**
	 * 'Send article to be transferred' (RFC977)
	 */
	ResponseCodeTransferSend = 335

	/**
	 * 'Article transferred ok' (RFC977)
	 */
	ResponseCodeTransferSuccess = 235

	/**
	 * 'Article not wanted - do not send it' (RFC977)
	 */
	ResponseCodeTransferUnwanted = 435

	/**
	 * 'Transfer failed - try again later' (RFC977)
	 */
	ResponseCodeTransferFailure = 436

	/**
	 * 'Article rejected - do not try again' (RFC977)
	 */
	ResponseCodeTransferRejected = 437

	// Posting

	/**
	 * 'Send article to be posted' (RFC977)
	 */
	ResponseCodePostingSend = 340

	/**
	 * 'Article posted ok' (RFC977)
	 */
	ResponseCodePostingSuccess = 240

	/**
	 * 'Posting not allowed' (RFC977)
	 */
	ResponseCodePostingProhibited = 440

	/**
	 * 'Posting failed' (RFC977)
	 */
	ResponseCodePostingFailure = 441

	// Authorization

	/**
	 * 'Authorization required for this command' (RFC2980)
	 */
	ResponseCodeAuthorizationRequired = 450

	/**
	 * 'Continue with authorization sequence' (RFC2980)
	 */
	ResponseCodeAuthorizationContinue = 350

	/**
	 * 'Authorization accepted' (RFC2980)
	 */
	ResponseCodeAuthorizationAccepted = 250

	/**
	 * 'Authorization rejected' (RFC2980)
	 */
	ResponseCodeAuthorizationRejected = 452

	// Authentication

	/**
	 * 'Authentication required' (RFC2980)
	 */
	ResponseCodeAuthenticationRequired = 480

	/**
	 * 'More authentication information required' (RFC2980)
	 */
	ResponseCodeAuthenticationContinue = 381

	/**
	 * 'Authentication accepted' (RFC2980)
	 */
	ResponseCodeAuthenticationAccepted = 281

	/**
	 * 'Authentication rejected' (RFC2980)
	 */
	ResponseCodeAuthenticationRejected = 482

	// Misc

	/**
	 * 'Help text follows' (Draft)
	 */
	ResponseCodeHelpFollows = 100

	/**
	 * 'Capabilities list follows' (Draft)
	 */
	ResponseCodeCapabilitiesFollow = 101

	/**
	 * 'Server date and time' (Draft)
	 */
	ResponseCodeServerDate = 111

	/**
	 * 'Information follows' (Draft)
	 */
	ResponseCodeInformationFollows = 215

	/**
	 * 'Groups follows' (Draft)
	 */
	ResponseCodeGroupsFollow = 215

	/**
	 * 'Overview information follows' (Draft)
	 */
	ResponseCodeOverviewFollows = 224

	/**
	 * 'Headers follow' (Draft)
	 */
	ResponseCodeHeadersFollow = 225

	/**
	 * 'List of new articles follows' (Draft)
	 */
	ResponseCodeNewArticlesFollow = 230

	/**
	 * 'List of new newsgroups follows' (Draft)
	 */
	ResponseCodeNewGroupsFollow = 231

	/**
	 * 'The server is in the wrong mode; the indicated capability should be used to change the mode' (Draft)
	 */
	ResponseCodeWrongMode = 401

	/**
	 * 'Internal fault or problem preventing action being taken' (Draft)
	 */
	ResponseCodeInternalFault = 403

	/**
	 * 'Command unavailable until suitable privacy has been arranged' (Draft)
	 *
	 * (the client must negotiate appropriate privacy protection on the connection.
	 * This will involve the use of a privacy extension such as [NNTP-TLS].)
	 */
	ResponseCodeEncryptionRequired = 483

	/**
	 * 'Error in base64-encoding [RFC3548] of an argument' (Draft)
	 */
	ResponseCodeBase64EncodingError = 504
)

type GroupPermission byte

const (
	GroupPostingPermitted GroupPermission = 'y'
	GroupPostingForbidden GroupPermission = 'n'
	GroupPostingModerated GroupPermission = 'm'
)

const DefaultArticleDateLayout = "02 Jan 2006 15:04:05 -0700"

type OverviewFieldType int

const (
	ShortHeaderOverviewField OverviewFieldType = iota
	FullHeaderOverviewField
	MetadataOverviewField
)
