import _ from 'lodash';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import Paper from '@material-ui/core/Paper';
import PropTypes from 'prop-types';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import { withStyles } from '@material-ui/core/styles';

const styles = theme => ({
  root: {
    width: '100%',
    marginTop: theme.spacing.unit * 3,
    overflowX: 'auto',
  },
  table: {
    minWidth: 700
  },
});

class BaseTable extends React.Component {
  handleRowExpand() {
    console.log("expanded");
  }

  render() {
    const { classes, tableRows, tableColumns, tableClassName } = this.props;
    return (
      <Paper className={classes.root}>
        <Table className={`${classes.table} ${tableClassName}`}>
          <TableHead>
            <TableRow>
              { _.map(tableColumns, c => (
                <TableCell
                  key={c.key}
                  numeric={c.isNumeric}>{c.title}
                </TableCell>
            ))
            }
            </TableRow>
          </TableHead>
          <TableBody>
            {
            _.map(tableRows, d => {
            return (
              <ExpansionPanel expanded={true} onChange={this.handleRowExpand('panel1')}>
                <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
                  <TableRow key={d.key}>
                    { _.map(tableColumns, c => (
                      <TableCell
                        key={`table-${d.key}-${c.key}`}
                        numeric={c.isNumeric}>{c.render(d)}
                      </TableCell>
                  ))
                }
                  </TableRow>
                </ExpansionPanelSummary>
                <ExpansionPanelDetails>
                  <div>
                  Nulla facilisi. Phasellus sollicitudin nulla et quam mattis feugiat. Aliquam eget
                  maximus est, id dignissim quam.
                  </div>
                </ExpansionPanelDetails>
              </ExpansionPanel>
            );
          })}
          </TableBody>
        </Table>
      </Paper>
    );
  }
}

BaseTable.propTypes = {
  classes: PropTypes.shape({}).isRequired,
  tableClassName: PropTypes.string,
  tableColumns: PropTypes.arrayOf(PropTypes.shape({
    title: PropTypes.string,
    isNumeric: PropTypes.bool,
    render: PropTypes.func
  })).isRequired,
  tableRows: PropTypes.arrayOf(PropTypes.shape({}))
};

BaseTable.defaultProps = {
  tableClassName: "",
  tableRows: []
};

export default withStyles(styles)(BaseTable);