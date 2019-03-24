package main

import (
	"context"
	"fmt"
	"os"
	"bufio"
	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
)

// function to sort releases array in descending order
func SortVersions(releases []*semver.Version) []*semver.Version {
  for i := 0; i < len(releases); i++ {
    for j := 0; j < len(releases) - 1 - i; j++ {
      if (releases[j].LessThan(*releases[j + 1])) {
        releases[j], releases[j + 1] = releases[j + 1], releases[j]
      }
    }
  }
  return releases
}

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	var versionSlice []*semver.Version
  releases = SortVersions(releases)
  if (minVersion.LessThan(*releases[0]) || minVersion.Equal(*releases[0])){
    versionSlice = append(versionSlice, releases[0])
  }
  index_version := 0
  if (len(versionSlice) > 0){
    for i := 1; i < len(releases); i++ {
      if (releases[i].Major != versionSlice[index_version].Major && (minVersion.LessThan(*releases[i]) || minVersion.Equal(*releases[i]))){
        versionSlice = append(versionSlice, releases[i])
        index_version ++
      }
      if (releases[i].Major == versionSlice[index_version].Major && releases[i].Minor != versionSlice[index_version].Minor && (minVersion.LessThan(*releases[i]) || minVersion.Equal(*releases[i]))){
        versionSlice = append(versionSlice, releases[i])
        index_version ++
      }
    }
  }
	return versionSlice
}

// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func main() {
	// variable declarations
	line := ""
	counter1 := 0
	counter2 := 0
	// Github
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 10}
	// read input file
	myfile, err := os.Open(os.Args[1])
	if err != nil {
	  fmt.Printf("Error opening input file\n")
	  os.Exit(1)
	  }
	  defer myfile.Close()
	  scanner := bufio.NewScanner(myfile)
	 // reading input file
	  for scanner.Scan() {
				first_arg := ""
				second_arg := ""
				version := ""
	      line = scanner.Text();
	  // parsing line into separate strings
	      for i:= 0; i < len(line); i++ {
	        if (line[i] == '/') {
	          counter1 = i;
	          break
	        } else {first_arg = first_arg + string(line[i])}
	      }
	      for j:= counter1 + 1; j < len(line); j++ {
	        if (line[j] == ',') {
	          counter2 = j;
	          break
	        } else {second_arg = second_arg + string(line[j])}
	      }
	      for k:= counter2 + 1; k < len(line); k++ {
	        version = version + string(line[k])
	      }
				releases, _, err := client.Repositories.ListReleases(ctx, first_arg, second_arg, opt)
				if err != nil {
					fmt.Println(err)
				}
				minVersion := semver.New(version)
				allReleases := make([]*semver.Version, len(releases))
				for i, release := range releases {
					versionString := *release.TagName
					if versionString[0] == 'v' {
						versionString = versionString[1:]
					}
					allReleases[i] = semver.New(versionString)
				}
				versionSlice := LatestVersions(allReleases, minVersion)
			// printing output
				fmt.Printf("latest versions of %s", first_arg)
				fmt.Printf("/%s", second_arg)
				fmt.Printf(": %s", versionSlice)
				fmt.Printf("\n")
	    }
	    if err := scanner.Err(); err != nil {
	        fmt.Printf("Invalid: %s ", err)
	    }
} // end of main
