####################################################################################################################################
# PROTOCOL COMMON MINION MODULE
####################################################################################################################################
package BackRest::Protocol::CommonMinion;
use parent 'BackRest::Protocol::Common';

use strict;
use warnings FATAL => qw(all);
use Carp qw(confess);

use File::Basename qw(dirname);
use Scalar::Util qw(blessed);

use lib dirname($0) . '/../lib';
use BackRest::Common::Exception;
use BackRest::Common::Ini;
use BackRest::Common::Log;
use BackRest::Common::String;
use BackRest::Config::Config;
use BackRest::Protocol::Common;
use BackRest::Protocol::IO;

####################################################################################################################################
# Operation constants
####################################################################################################################################
use constant OP_PROTOCOL_COMMON_MINION                              => 'Protocol::CommonMinion';

use constant OP_PROTOCOL_COMMON_MINION_NEW                          => OP_PROTOCOL_COMMON_MINION . "->new";

####################################################################################################################################
# CONSTRUCTOR
####################################################################################################################################
sub new
{
    my $class = shift;                  # Class name

    # Assign function parameters, defaults, and log debug info
    my
    (
        $strOperation,
        $strName,                                   # Name of the protocol
        $iBlockSize,                                # Buffer size
        $iCompressLevel,                            # Set compression level
        $iCompressLevelNetwork                      # Set compression level for network only compression
    ) =
        logDebugParam
        (
            OP_PROTOCOL_COMMON_MINION_NEW, \@_,
            {name => 'strName'},
            {name => 'iBlockSize'},
            {name => 'iCompressLevel'},
            {name => 'iCompressLevelNetwork'}
        );

    # Create the class hash
    my $self = $class->SUPER::new($iBlockSize, $iCompressLevel, $iCompressLevelNetwork, $strName);
    bless $self, $class;

    # Create the IO object with std io
    $self->{io} = new BackRest::Protocol::IO(*STDIN, *STDOUT, *STDERR);

    # Write the greeting so master process knows who we are
    $self->greetingWrite();

    return $self;
}

####################################################################################################################################
# greetingWrite
#
# Send a greeting to the master process.
####################################################################################################################################
sub greetingWrite
{
    my $self = shift;

    $self->{io}->lineWrite($self->{strGreeting});
}

####################################################################################################################################
# textWrite
#
# Write text out to the protocol layer prefixing each line with a period.
####################################################################################################################################
sub textWrite
{
    my $self = shift;
    my $strBuffer = shift;

    $strBuffer =~ s/\n/\n\./g;

    $self->{io}->lineWrite('.' . $strBuffer);
}

####################################################################################################################################
# binaryXferAbort
#
# Abort transfer when source file does not exist.
####################################################################################################################################
sub binaryXferAbort
{
    my $self = shift;

    # Only allow in the backend process
    $self->{io}->lineWrite('block -1');
}

####################################################################################################################################
# errorWrite
#
# Write errors with error codes in protocol format, otherwise write to stderr and exit with error.
####################################################################################################################################
sub errorWrite
{
    my $self = shift;
    my $oMessage = shift;

    my $iCode;
    my $strMessage;

    # If the message is blessed it may be a standard exception
    if (blessed($oMessage))
    {
        # Check if it is a standard exception
        if ($oMessage->isa('BackRest::Common::Exception'))
        {
            $iCode = $oMessage->code();
            $strMessage = $oMessage->message();
        }
        # Else terminate the process with an error
        else
        {
            $self->{io}->errorWrite('unknown error object');
            exit ERROR_UNKNOWN;
        }
    }
    # Else terminate the process with an error
    else
    {
        $self->{io}->errorWrite('unknown error: ' . $oMessage);
        exit ERROR_UNKNOWN;
    }

    # Write the message text into protocol
    if (defined($strMessage))
    {
        $self->textWrite(trim($strMessage));
    }

    # Indicate that an error ocurred and provide the code
    $self->{io}->lineWrite("ERROR" . (defined($iCode) ? " $iCode" : ''));
}

####################################################################################################################################
# outputWrite
#
# Write output for the master process.
####################################################################################################################################
sub outputWrite
{
    my $self = shift;
    my $strOutput = shift;

    if (defined($strOutput))
    {
        $self->textWrite($strOutput);
    }

    $self->{io}->lineWrite('OK');
}

####################################################################################################################################
# cmdRead
#
# Read command sent by the master process.
####################################################################################################################################
sub cmdRead
{
    my $self = shift;
    my $oParamHashRef = shift;

    my $strLine;
    my $strCommand;

    while ($strLine = $self->{io}->lineRead())
    {
        if (!defined($strCommand))
        {
            if ($strLine =~ /:$/)
            {
                $strCommand = substr($strLine, 0, length($strLine) - 1);
            }
            else
            {
                $strCommand = $strLine;
                last;
            }
        }
        else
        {
            if ($strLine eq 'end')
            {
                last;
            }

            my $iPos = index($strLine, '=');

            if ($iPos == -1)
            {
                confess "param \"${strLine}\" is missing = character";
            }

            my $strParam = substr($strLine, 0, $iPos);
            my $strValue = substr($strLine, $iPos + 1);

            ${$oParamHashRef}{"${strParam}"} = ${strValue};
        }
    }

    return $strCommand;
}

1;
