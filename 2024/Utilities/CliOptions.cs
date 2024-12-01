using CommandLine;

namespace AdventOfCode.Utilities;

public class CliOptions
{
    [Option('d', Required = false, Default = -1, HelpText = "Select which day to execute")]
    public int Day { get; set; }
    [Option('p', Required = false, Default = -1, HelpText = "Select which part of the day to execute")]
    public int Part { get; set; }
    [Option('t', Required = false, Default = false, HelpText = "Flag for running the part's input test cases")]
    public bool RunTest { get; set; }
    [Value(0, MetaName = "InputPath", Required = false, Default = "./input", HelpText = "Path to input files")]
    public string InputPath { get; set; }
}
